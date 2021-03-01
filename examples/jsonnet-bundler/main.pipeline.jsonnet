local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';

function(params={}) (
	pipeline.Pipeline {
		name: "name",
		application: "application",
	}
)