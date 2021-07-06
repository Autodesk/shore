/**
    Creates a pipeline.
**/

function(params={}) (
	{
		"application": params["application"],
		"pipeline": params["pipeline"],
		"message": "Hello %s!" % [params["example_value"]],
	}
)
