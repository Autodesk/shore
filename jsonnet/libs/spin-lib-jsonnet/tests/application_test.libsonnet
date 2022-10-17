local application = import '../application.libsonnet';

// Data to be used for tests
local banner = application.Banner {
  text: 'my banner',
};

local dataSourcesConfig = application.DataSourcesConfiguration {
  enabled: [
    application.DataSource.serverGroups,
    application.DataSource.loadBalancers,
    application.DataSource.executions,
  ],
  disabled: [
    application.DataSource.securityGroups,
    application.DataSource.functions,
  ],
};

local tests = [
  application.Application {
    cloudProviders: 'aws, kubernetes',
    email: 'email',
    description: 'description',
    name: 'name',
  },
  banner,
  dataSourcesConfig,
  application.NewApplication([banner], dataSourcesConfig) {
    cloudProviders: 'aws, kubernetes',
    email: 'email',
    description: 'description',
    name: 'name',
  },
  application.NewApplication([], {}) {
    cloudProviders: 'aws, kubernetes',
    email: 'email',
    description: 'description',
    name: 'name',
    providerSettings: application.AwsProviderSettings,
  },
];

local assertions = [
  {
    cloudProviders: 'aws, kubernetes',
    description: 'description',
    email: 'email',
    name: 'name',
    platformHealthOnly: false,
    providerSettings: {},
  },
  {
    backgroundColor: 'var(--color-alert)',
    enabled: true,
    text: 'my banner',
    textColor: 'var(--color-text-on-dark)',
  },
  {
    disabled: [
      'securityGroups',
      'functions',
    ],
    enabled: [
      'serverGroups',
      'loadBalancers',
      'executions',
    ],
  },
  {
    cloudProviders: 'aws, kubernetes',
    customBanners: [
      {
        backgroundColor: 'var(--color-alert)',
        enabled: true,
        text: 'my banner',
        textColor: 'var(--color-text-on-dark)',
      },
    ],
    dataSources: {
      disabled: [
        'securityGroups',
        'functions',
      ],
      enabled: [
        'serverGroups',
        'loadBalancers',
        'executions',
      ],
    },
    description: 'description',
    email: 'email',
    name: 'name',
    platformHealthOnly: false,
    providerSettings: {},
  },
  {
    cloudProviders: 'aws, kubernetes',
    description: 'description',
    email: 'email',
    name: 'name',
    platformHealthOnly: false,
    providerSettings: {
      aws: {
        useAmiBlockDeviceMappings: false,
      },
    },
  },
];

{
  pass: tests == assertions,
  tests: tests,
  assertions: assertions,
}
