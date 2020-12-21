#include <assert.h>
#include <signal.h>
#include <stdio.h>
#include <pulse/pulseaudio.h>

pa_mainloop* _mainloop;
pa_mainloop_api* _mainloop_api;
pa_context* _context;
pa_signal_event* _signal;
const char* default_sink_name;

extern void goSetVolume(float);


/*
* Called when the requested sink information is ready.
*/
static void sink_info_callback(pa_context *c, const pa_sink_info *i,
		int eol, void *userdata)
{
	if (i)
	{
		float volume = (float)pa_cvolume_avg(&(i->volume)) / (float)PA_VOLUME_NORM;
		if (i->mute) {
			volume = 0;
		}
		goSetVolume(volume);
	}
}


/*
* Called when the requested information on the server is ready. This is
* used to find the default PulseAudio sink.
*/
static void server_info_callback(pa_context *c, const pa_server_info *i,
		void *userdata)
{
	default_sink_name = i->default_sink_name;
	pa_context_get_sink_info_by_name(c, i->default_sink_name, sink_info_callback, userdata);
}


/*
	* Called when an event we subscribed to occurs.
	*/
static void subscribe_callback(pa_context *c,
		pa_subscription_event_type_t type, uint32_t idx, void *userdata)
{
	pa_context_get_server_info(c, server_info_callback, userdata);
}

/*
	* Called whenever the context status changes.
	*/
static void context_state_callback(pa_context *c, void *userdata)
{
	switch (pa_context_get_state(c))
	{
		case PA_CONTEXT_CONNECTING:
		case PA_CONTEXT_AUTHORIZING:
		case PA_CONTEXT_SETTING_NAME:
			break;

		case PA_CONTEXT_READY:
			pa_context_get_server_info(c, server_info_callback, userdata);

			// Subscribe to sink events from the server. This is how we get
			// volume change notifications from the server.
			pa_context_set_subscribe_callback(c, subscribe_callback, userdata);
			pa_context_subscribe(
				c,
				PA_SUBSCRIPTION_MASK_SINK|PA_SUBSCRIPTION_MASK_SINK_INPUT|PA_SUBSCRIPTION_MASK_SERVER,
				NULL,
				NULL
			);
			break;
	}
}


/**
 * Initializes state and connects to the PulseAudio server.
 */
int initialize()
{
	_mainloop = pa_mainloop_new();
	if (!_mainloop)
	{
		fprintf(stderr, "pa_mainloop_new() failed.\n");
		return 1;
	}

	_mainloop_api = pa_mainloop_get_api(_mainloop);

	if (pa_signal_init(_mainloop_api) != 0)
	{
		fprintf(stderr, "pa_signal_init() failed\n");
		return 1;
	}


	_context = pa_context_new(_mainloop_api, "PulseAudio Test");
	if (!_context)
	{
		fprintf(stderr, "pa_context_new() failed\n");
		return 1;
	}

	pa_context_set_state_callback(_context, context_state_callback, NULL);

	if (pa_context_connect(_context, NULL, PA_CONTEXT_NOAUTOSPAWN, NULL) < 0)
	{
		fprintf(stderr, "pa_context_connect() failed: %s\n", pa_strerror(pa_context_errno(_context)));
		return 1;
	}

	return 0;
}

/**
 * Runs the main PulseAudio event loop. Calling quit will cause the event
 * loop to exit.
 */
int run()
{
	int ret = 1;
	if (pa_mainloop_run(_mainloop, &ret) < 0)
	{
		fprintf(stderr, "pa_mainloop_run() failed.\n");
		return ret;
	}

	return ret;
}

/**
 * Exits the main loop with the specified return code.
 */
void quit(int ret)
{
	_mainloop_api->quit(_mainloop_api, ret);
}

/**
 * Called when the PulseAudio system is to be destroyed.
 */
void destroy()
{
	if (_context)
	{
		pa_context_unref(_context);
		_context = NULL;
	}

	if (_signal)
	{
		pa_signal_free(_signal);
		pa_signal_done();
		_signal = NULL;
	}

	if (_mainloop)
	{
		pa_mainloop_free(_mainloop);
		_mainloop = NULL;
		_mainloop_api = NULL;
	}
}
