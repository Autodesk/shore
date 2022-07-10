local artifact = import 'spin-lib-jsonnet/artifact.libsonnet';
local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';


function(params) (
  pipeline.Pipeline {
    application: params.application,
    name: params.pipeline,
    parameterConfig: [
      parameter.Parameter {
        name: 'kube-deployment',
        required: true,
      },
      parameter.Parameter {
        name: 'kube-service',
        required: true,
      },
    ],
    Stages:: [
      {
        account: 'kubernetes',
        alias: 'deployManifest',
        application: '',
        cloudProvider: 'kubernetes',
        consumeArtifactSource: 'none',
        manifests: [
          '${#readJson(parameters["kube-deployment"])}',
          '${#readJson(parameters["kube-deployment"])}',
        ],
        manifestArtifact: null,
        manifestArtifactId: null,
        moniker: {
          app: 'test1test2test3',
        },
        name: 'deploy manifest',
        skipExpressionEvaluation: false,
        source: 'text',
        trafficManagement: {
          enabled: false,
          options: {
            enableTraffic: false,
            namespace: null,
            services: null,
          },
        },
        type: 'deployManifest',
      },
    ],
  }
)
