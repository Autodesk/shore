/**
   @file Contains objects for Kubernetes Job configuration.
**/

/**
  Creates a base Job-Spec/Manifest.

  Used primarily for extending as Manifest.

  @example
    local myContainer = kube.Container { ... };

    ...

    kube.JobSpec {
        metadata: {
            name: "my-kube-job",
            labels: { ... }
        },
        spec: {
            template: { ... }
        }
    }

    @property {Object} metadata - Kube metadata to use.
    @property {Object} spec - Kube spec to use.
**/
local JobSpec = {
  metadata: error '`metadata` (Object) property is required for JobSpec',
  spec: error '`spec` (Object) property is required for JobSpec',

  apiVersion: 'batch/v1',
  kind: 'Job',
};

/**
  Creates a Container object that's used in the Manifest.

  @example
    kubeManifest.container {
        name: 'potato-container',
        image: 'localhost/potato-image:latest',
    }

    @property {String} image - Image for the container.
    @property {String} name - Name of the container.
    @property {Array<String>} [args=[]] - The container arguments.
    @property {Array<String>} [command=[]] - The container command.
    @property {Array<String>} [env=[]] - The container environment variables.
    @property {Array<VolumeMounts>} [volumeMounts=[]] - The associated volume mounts in the container.
**/
local Container = {
  image: error '`image` (String) property is required for Container',
  name: error '`name` (String) property is required for Container',

  args: [],
  command: [],
  env: [],
  volumeMounts: [],
};

/**
  Creates a SecretVolume object that's used in the Manifest.

  @example

    vol1 = kube.SecretVolume { ... };

    vol2 = kube.SecretVolume {
      name: "secvolume",
      secretName: "secname",
    };

    kube.Manifest = {
      generateName: 'tomato',
      namespace:: 'potato'
      containers: [..]
      volumes: [vol1, vol2]
    };

  @property {String} name - Name of the SecretVolume.
  @property {String} secretName - The secret to mount for the SecretVolume.

**/
local SecretVolume = {
  local this = self,
  name: error '`name` (String) property is required for SecretVolume',
  secretName:: error '`secretName` (String) property is required for SecretVolume',

  secret: {
    secretName: this.secretName,
  },
};

/**

  Creates a VolumeMount object that's used in the Manifest.

  @example

    vm1 = kube.VolumeMount { ... };

    vm2 = kube.VolumeMount {
        mountPath: "/etc/secrets"
        name: "secvolume",
        readOnly: true,
    };

    kubeManifest.container {
        name: 'potato-container',
        image: 'localhost/potato-image:latest',
        volumeMounts: [vm1, vm2],
    }

  @property {String} mountPath - Location to mount the volume in the container.
  @property {String} name - Name of the VolumeMount. Corresponds to the SecretVolume name in the container definition.
  @property {Bool} readOnly - the r/w attribute of the volumeMount.

**/
local VolumeMount = {
  mountPath: error '`mountPath` (String) property is required for VolumeMount',
  name: error '`name` (String) property is required for VolumeMount',
  readOnly: true,
};

/**
  Creates a Kubernetes Manifest.

  This can be used for the RunKubeJobStage in the stage.libsonnet.

  It is encouraged to use the `labels` property to have some organization for the containers spun up on the backend.

  `name` vs `generatedName`:

  `name` needs to unique in the namespace that the container is running.

  `generatedName` will be suffixed with a random hex string to make a unique `name`.

  More info:

  https://kubernetes.io/docs/reference/using-api/api-concepts/

  https://serverfault.com/questions/809632/is-it-possible-to-rerun-kubernetes-job/868826#868826

  @example
    local myContainer = kube.Container { ... };

    ...

    kube.Manifest {
        generateName: 'my-kube-job',
        containers: [
            myContainer
        ],

        labels:: {
          purpose: 'spinnaker-run-kube-job-stage',
          spinnakerApplication: spel.makeSanitizedPipelineExecutionAppSpelExp,
          spinnakerPipeline: spel.makeSanitizedPipelineExecutionNameSpelExp,
          pipelineExecId: spel.makePipelineExecutionIdSpelExp,
        }
    }

    @class
    @augments JobSpec

    @property {String} generateName - generateName of the Kube Run Job.
    @property {Array<Container>} containers - An array of containers this job wills run.
    @property {int} [backoffLimit=6] - The number of times to backoff.
    @property {String} [namespace='default'] - The namespace of the Kube Run Job.
    @property {String} [restartPolicy='Never'] - Whether or not to restart the pods.
    @property {Array<SecretVolume>} [volumes=[]] volumes - An arry of the SecretVolumes spec to mount secrets for containers in Run Job defintions.
    @property {String} [serviceAccoutName='default'] - The serviceAccountName user the Run Job will run with.
    @property {int} [ttlSecondsAfterFinished=null] - Cleans up the job after the specified amount of seconds passed after its' completion. Encouraged to be set to keep the kube cluster clean and healthy.
**/
local Manifest = JobSpec {
  local this = self,

  generateName:: error '`generateName` (String) property is required for Container',
  containers:: error '`containers` (Array<Container>) property is required for Container',

  backoffLimit:: 6,  // 6 is the default according to kube docs
  restartPolicy:: 'Never',
  namespace:: 'default',
  labels:: {},
  serviceAccountName:: 'default',
  volumes:: [],
  ttlSecondsAfterFinished:: null,

  metadata: {
    generateName: this.generateName,
    labels: this.labels,
    namespace: this.namespace,
  },
  spec: {
    backoffLimit: this.backoffLimit,
    [if this.ttlSecondsAfterFinished != null then 'ttlSecondsAfterFinished']: this.ttlSecondsAfterFinished,
    template: {
      spec: {
        serviceAccountName: this.serviceAccountName,
        volumes: this.volumes,
        restartPolicy: this.restartPolicy,
        containers: this.containers,
      },
    },
  },
};

// Exposed for public use.
{
  JobSpec:: JobSpec,

  Container:: Container,
  SecretVolume:: SecretVolume,
  VolumeMount:: VolumeMount,
  Manifest:: Manifest,
}
