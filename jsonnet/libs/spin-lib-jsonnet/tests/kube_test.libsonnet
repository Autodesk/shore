local kube = import '../kube.libsonnet';

local tests = [
  kube.JobSpec {
    metadata: {},
    spec: {},
  },
  kube.Container {
    name: 'potato-container',
    image: 'localhost/potato-image:latest',
    command: ['./fry'],
    args: ['--with-salt'],
    volumeMounts: [kube.VolumeMount { mountPath: '/tmp/secrets', name: 'volname' }],
  },

  kube.Manifest {
    generateName:: 'my-kube-job',
    containers:: [
      kube.Container { name: 'potato-container', image: 'localhost/potato-image:latest' },
    ],
    labels:: {
      purpose: 'spinnaker-run-kube-job-stage',
      spinnakerApplication: 'my-spinnaker-app',
      spinnakerPipeline: 'my-amazing-pipeline',
    },
    namespace:: 'my-namespace',
  },
  kube.Manifest {
    generateName:: 'my-kube-job',
    serviceAccountName:: 'my-svc',
    volumes:: [kube.SecretVolume { name: 'volname', secretName:: 'secname' }],
    containers:: [
      kube.Container {
        name: 'potato-container',
        image: 'localhost/potato-image:latest',
        volumeMounts: [{
          mountPath: '/etc/mysecret',
          name: 'volname',
          readOnly: true,
        }],
      },
    ],
  },
  kube.Manifest {
    generateName:: 'my-kube-job',
    ttlSecondsAfterFinished:: 12345,
    containers:: [
      kube.Container { name: 'potato-container', image: 'localhost/potato-image:latest' },
    ],
    labels:: {
      purpose: 'spinnaker-run-kube-job-stage',
      spinnakerApplication: 'my-spinnaker-app',
      spinnakerPipeline: 'my-amazing-pipeline',
    },
    namespace:: 'my-namespace',
  },
];

local assertions = [
  {
    apiVersion: 'batch/v1',
    kind: 'Job',
    metadata: {},
    spec: {},
  },
  {
    args: [
      '--with-salt',
    ],
    command: [
      './fry',
    ],
    env: [],
    image: 'localhost/potato-image:latest',
    name: 'potato-container',
    volumeMounts: [{
      mountPath: '/tmp/secrets',
      name: 'volname',
      readOnly: true,
    }],
  },
  {
    apiVersion: 'batch/v1',
    kind: 'Job',
    metadata: {
      labels: {
        purpose: 'spinnaker-run-kube-job-stage',
        spinnakerApplication: 'my-spinnaker-app',
        spinnakerPipeline: 'my-amazing-pipeline',
      },
      generateName: 'my-kube-job',
      namespace: 'my-namespace',
    },
    spec: {
      backoffLimit: 6,
      template: {
        spec: {
          containers: [
            {
              args: [],
              command: [],
              env: [],
              image: 'localhost/potato-image:latest',
              name: 'potato-container',
              volumeMounts: [],
            },
          ],
          restartPolicy: 'Never',
          serviceAccountName: 'default',
          volumes: [],
        },
      },
    },
  },
  {
    apiVersion: 'batch/v1',
    kind: 'Job',
    metadata: {
      generateName: 'my-kube-job',
      labels: {},
      namespace: 'default',
    },
    spec: {
      backoffLimit: 6,
      template: {
        spec: {
          containers: [
            {
              args: [],
              command: [],
              env: [],
              name: 'potato-container',
              image: 'localhost/potato-image:latest',
              volumeMounts: [{
                mountPath: '/etc/mysecret',
                name: 'volname',
                readOnly: true,
              }],
            },
          ],
          restartPolicy: 'Never',
          serviceAccountName: 'my-svc',
          volumes: [
            {
              name: 'volname',
              secret: {
                secretName: 'secname',
              },
            },
          ],
        },
      },
    },
  },
  {
    apiVersion: 'batch/v1',
    kind: 'Job',
    metadata: {
      generateName: 'my-kube-job',
      labels: {
        purpose: 'spinnaker-run-kube-job-stage',
        spinnakerApplication: 'my-spinnaker-app',
        spinnakerPipeline: 'my-amazing-pipeline',
      },
      namespace: 'my-namespace',
    },
    spec: {
      backoffLimit: 6,
      template: {
        spec: {
          containers: [
            {
              args: [],
              command: [],
              env: [],
              image: 'localhost/potato-image:latest',
              name: 'potato-container',
              volumeMounts: [],
            },
          ],
          restartPolicy: 'Never',
          serviceAccountName: 'default',
          volumes: [],
        },
      },
      ttlSecondsAfterFinished: 12345,
    },
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
