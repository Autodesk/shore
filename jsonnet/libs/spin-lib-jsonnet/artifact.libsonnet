/**
    @file Contains objects for creating various Spinnaker Artifacts.

    More information on what is a Spinnaker Artifact can be found at:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/#what-is-a-spinnaker-artifact}
**/

/**
    A generic base Artifact object.

    Used for extension by other  Artifact object types.


    A list of all Artifact types can be found in OSS deck's repository:

    {@link https://github.com/spinnaker/deck/blob/master/app/scripts/modules/core/src/artifact/ArtifactTypes.ts}


    Or in the ArtifactCredential classes of clouddriver:

    {@link https://github.com/spinnaker/clouddriver/tree/master/clouddriver-artifacts/src/main/java/com/netflix/spinnaker/clouddriver/artifacts}

    @example
        local myArtifact = artifact.Artifact {
            artifactAccount: 'custom-artifact',
            name: 'myArtifact',
            type: 'custom/object',
            reference: 'Some say tomatoes, some say to-mah-toes.',
        };

    @property {String} artifactAccount - The Spinnaker artifact account.
    @property {String} name - Name of the artifact. Used as a handle throughout the pipeline.
    @property {String} type - The type of artifact.
    @property {String} reference - Reference/Value of the artifact.
    @property {String} [displayName] - The name of the artifact displayed in the UI. Has no impact on usage.
**/
local Artifact = {
  artifactAccount: error '`artifactAccount` (String) property is required for Artifact',
  name: error '`name` (String) property is required for Artifact',
  type: error '`type` (String) property is required for Artifact',
  reference: error '`reference` (String) property is required for Artifact',

  // Other properties
  // displayName: '',
};


/**
    Custom Artifact.

    General-purpose artifact type. Can be used to hold any kind of string as a refernce.

    Code:

    {@link https://github.com/spinnaker/deck/blob/master/app/scripts/modules/core/src/pipeline/config/triggers/artifacts/custom/CustomArtifactEditor.tsx}

    @example
        local myArtifact = artifact.CustomArtifact {
            name: 'myArtifact',
            reference: 'Some say tomatoes, some say to-mah-toes.',
        };

    @class
    @augments Artifact
**/
local CustomArtifact = {
  artifactAccount: 'custom-artifact',
  type: 'custom/object',
};

/**
    GitHub File Artifact.

    Full documentation found here:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/types/github-file/}

    @example
        artifact.GitHubFileArtifact {
            artifactAccount: 'spinnaker-github-account',
            name: 'service-manifest'
            reference: 'https://<GitHub URL>/service.manifest',
            version: 'my-dev-branch'
        }
    @class
    @augments Artifact

    @property {String} [location=''] - The directory inside of the project, if not at root.
**/
local GitHubFileArtifact = Artifact {
  location: '',
  type: 'github/file',
  version: 'master',
};

/**
    Git Repo Artifact.

    Full documentation found here:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/types/git-repo/}

    @example
        artifact.GitRepoArtifact {
            artifactAccount: 'spinnaker-gitrepo-account',
            reference: 'https://<GitHub URL>/service-repo.git',
            name: 'service-repo',
        }

    @class
    @augments Artifact
**/
local GitRepoArtifact = Artifact {
  type: 'git/repo',
  version: 'master',
};

/**
    HTTP File Artifact.

    Full documentation found here:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/types/http-file/}

    @example
        artifact.HttpFileArtifact {
            artifactAccount: 'spinnaker-http-account',
            reference: 'https://localhost/README.md',
            name: 'my-readme',
        }

    @class
    @augments Artifact
**/
local HttpFileArtifact = Artifact {
  type: 'http/file',
};

/**
    Embedded Base64 Artifact.

    Full documentation found here:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/types/embedded-base64/}

    @example
        artifact.EmbeddedBase64Artifact {
            reference: 'SGVsbG8gV29ybGQh',
        }

    @class
    @augments Artifact
**/
local EmbeddedBase64Artifact = Artifact {
  type: 'embedded/base64',
};

/**
    S3 Object Artifact.

    Full documentation found here:

    {@link https://spinnaker.io/reference/artifacts-with-artifactsrewrite/types/s3-object/}

    @example
        artifact.S3ObjectArtifact {
            reference: 's3://myBucket/service.manifest',
            location: 'us-west-2'
        }

    @class
    @augments Artifact
**/
local S3ObjectArtifact = Artifact {
  location: error '`location` (String) property is required for S3ObjectArtifact',

  type: 's3/object',
};

/**
    Creates a Spinnaker object that specifies the expected Artifacts - on stages or pipelines.

    A default artifact can be provided by `defaultArtifact` property in conjuction with setting `useDefaultArtifact` to
    true.

    The `matchArtifact` property must be an object that has Artifact properties that should be matched on. This can be a subset of the Artifact properties, whichever ones are relevant to match on.

    An artifact from the prior pipeline execution can be used by setting `usePriorArtifact` to true.

    More information on Expected Artifact can be found at:

    {@link https://spinnaker.io/reference/artifacts/in-pipelines/#expected-artifacts}

    @example
        local expectedGitRepoArtifact = {
            name: 'repository-name',
            type: 'git/repo'
            // Fields omitted by choice - either not important to match on, or a value that's expected to change.
            // version: '',
            // artifactAccount: '',
            // reference: '',
        }

        ...

        artifact.ExpectedArtifact {
            matchArtifact: expectedGitRepoArtifact
        }

    @property {Object} matchArtifact - This is an object that should have a subset of fields from Artifact object that is expected.
    @property {Artifact} [defaultArtifact={}] - The default artifact to provide if `useDefaultArtifact` is true.
    @property {displayName} [displayName] - The name of the artifact displayed in the UI. Does not effect functionality.
    @property {Boolean} [useDefaultArtifact=false] - Enables providing a default artifact when none match, if set to true.
    @property {Boolean} [usePriorArtifact=false] - When set to true and no artifact is matched, the artifact by the same name from the previous execution will be used. If no previous executions exist, the pipeline will error out.

    @class
    @augments Artifact
**/
local ExpectedArtifact = {
  matchArtifact: error '`matchArtifact` (Object) property is required for ExpectedArtifact',

  defaultArtifact: {},
  useDefaultArtifact: false,
  usePriorArtifact: false,

  // Other properties
  // displayName: '',
};

// Exposed for public use.
{
  Artifact:: Artifact,

  CustomArtifact:: CustomArtifact,
  EmbeddedBase64Artifact:: EmbeddedBase64Artifact,
  GitHubFileArtifact:: GitHubFileArtifact,
  GitRepoArtifact:: GitRepoArtifact,
  HttpFileArtifact:: HttpFileArtifact,
  S3ObjectArtifact:: S3ObjectArtifact,

  ExpectedArtifact:: ExpectedArtifact,
}
