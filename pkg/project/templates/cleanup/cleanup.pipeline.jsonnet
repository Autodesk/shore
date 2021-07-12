/**
    Creates a pipeline meant for cleanup.
**/

function(params={}) (
	{
		"application": params["application"],
		"pipeline": params["pipeline"],
		"message": "Hello %s! Please add some cleanup!" % [params["example_value"]],
	}
)
