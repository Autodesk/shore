local artifacts = import '../artifact.libsonnet';

// Data used for tests.
local httpFileArtifact = artifacts.HttpFileArtifact {
  name: 'name',
  artifactAccount: 'artifactAccount',
  reference: 'https://localhost/README.md',
};
local defaultHttpFileArtifact = artifacts.HttpFileArtifact {
  name: 'name',
  artifactAccount: 'artifactAccount',
  reference: 'https://localhost/DEFAULT.md',
};

local tests = [
  artifacts.Artifact {
    name: 'name',
    type: 'type',
    artifactAccount: 'artifactAccount',
    reference: 'reference',
  },
  artifacts.CustomArtifact {
    name: 'name',
    reference: 'potato',
  },
  artifacts.GitHubFileArtifact {
    name: 'name',
    artifactAccount: 'artifactAccount',
    reference: 'https://<GitHub URL>/service.manifest',
  },
  artifacts.GitRepoArtifact {
    name: 'name',
    artifactAccount: 'artifactAccount',
    reference: 'https://<GitHub URL>/service-repo.git',
  },
  httpFileArtifact,
  artifacts.EmbeddedBase64Artifact {
    name: 'name',
    artifactAccount: 'artifactAccount',
    reference: 'SGVsbG8gV29ybGQh',
  },
  artifacts.S3ObjectArtifact {
    name: 'name',
    artifactAccount: 'artifactAccount',
    reference: 's3://myBucket/service.manifest',
    location: 'us-west-2',
  },
  artifacts.ExpectedArtifact {
    matchArtifact: httpFileArtifact,
  },
  artifacts.ExpectedArtifact {
    matchArtifact: httpFileArtifact,
    useDefaultArtifact: true,
    defaultArtifact: defaultHttpFileArtifact,
  },
];

local assertions = [
  {
    artifactAccount: 'artifactAccount',
    name: 'name',
    type: 'type',
    reference: 'reference',
  },
  {
    artifactAccount: 'custom-artifact',
    name: 'name',
    type: 'custom/object',
    reference: 'potato',
  },
  {
    artifactAccount: 'artifactAccount',
    location: '',
    name: 'name',
    reference: 'https://<GitHub URL>/service.manifest',
    type: 'github/file',
    version: 'master',
  },
  {
    artifactAccount: 'artifactAccount',
    name: 'name',
    reference: 'https://<GitHub URL>/service-repo.git',
    type: 'git/repo',
    version: 'master',
  },
  {
    artifactAccount: 'artifactAccount',
    name: 'name',
    reference: 'https://localhost/README.md',
    type: 'http/file',
  },
  {
    artifactAccount: 'artifactAccount',
    name: 'name',
    reference: 'SGVsbG8gV29ybGQh',
    type: 'embedded/base64',
  },
  {
    artifactAccount: 'artifactAccount',
    location: 'us-west-2',
    name: 'name',
    reference: 's3://myBucket/service.manifest',
    type: 's3/object',
  },
  {
    defaultArtifact: {},
    matchArtifact: {
      artifactAccount: 'artifactAccount',
      name: 'name',
      reference: 'https://localhost/README.md',
      type: 'http/file',
    },
    useDefaultArtifact: false,
    usePriorArtifact: false,
  },
  {
    defaultArtifact: {
      artifactAccount: 'artifactAccount',
      name: 'name',
      reference: 'https://localhost/DEFAULT.md',
      type: 'http/file',
    },
    matchArtifact: {
      artifactAccount: 'artifactAccount',
      name: 'name',
      reference: 'https://localhost/README.md',
      type: 'http/file',
    },
    useDefaultArtifact: true,
    usePriorArtifact: false,
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
