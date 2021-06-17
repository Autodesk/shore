local artifact = import 'spin-lib-jsonnet/artifact.libsonnet';
local pipeline = import 'spin-lib-jsonnet/pipeline.libsonnet';
local stage = import 'spin-lib-jsonnet/stage.libsonnet';
local parameter = import 'spin-lib-jsonnet/parameter.libsonnet';

function(params) (
  pipeline.Pipeline {
    application: params.application,
    name: params.pipeline_name,
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
        cloudProvider: 'kubernetes',
        kinds: [
          'deployment',
          'service',
        ],
        labelSelectors: {
          selectors: [
            {
              key: 'app',
              kind: 'EQUALS',
              values: [
                '${#readJson(parameters["kube-deployment"]).metadata.labels.app}',
              ],
            },
          ],
        },
        location: 'default',
        mode: 'label',
        name: 'Delete (Manifest)',
        options: {
          cascading: true,
        },
        type: 'deleteManifest',
      },
    ],
  }
)
