/**
    @file Contains objects for creating a Spinnaker Application.
**/

/**
    Creates a Spinnaker Application object.

    Cloud providers can be found by looking for <implements CloudProvider> and their IDs in OSS clouddriver's repository:

    {@link https://github.com/spinnaker/clouddriver}

    Cloud providers must also be enabled for the Spinnaker environment in order for them to be usable.

    @example
        local banners =  [ application.Banner { ... }, ... ]

        ...

        application.Application {
            name: 'My Amazing Spinnaker App',
            description: 'Does amazing things with amazing services.',
            email: 'john.doe@shore.com',
            cloudProviders: 'aws',
            customBanners: banners,
        }

    @property {String} cloudProviders - The cloud providers that the Spinnaker Application will use.
    @property {String} email - The email of the owner of the Spinnaker Application.
    @property {String} description - A description about this Spinnaker Application.
    @property {String} name - The name of this Spinnaker Application.
    @property {Boolean} [platformHealthOnly=false] - Whether or not to check only the platform's health. For AWS EC2 instance, this means the instance is up, but might have passed healthchecks.
    @property {AwsProviderSettings} [providerSettings={}] - The provider settings for the application.
**/
local Application = {
  cloudProviders: error '`cloudProviders` (String) property is required for Application',
  email: error '`email` (String) property is required for Application',
  description: error '`description` (String) property is required for Application',
  name: error '`name` (String) property is required for Application',

  platformHealthOnly: false,
  providerSettings: {},
  // Other properties

  //customBanners: [],

  // Empty dataSources causing "Application Not Found" error in the web app? :thinking:
  //dataSources: {
  //  disable: [],
  //  enabled: [],
  //},
};

/**
    Creates a ProviderSettings object for AWS CloudProvider.

    This is specific to the AWS CloudProvider, and other providers will have their own settings.

    @example
        local myProviderSettings = application.AwsProviderSettings {
            useAmiBlockDeviceMappings:: true,
        }

        ...

        application.Application {
            ...
            providerSettings: myProviderSettings,
            ...
        }

    @property {Boolean} [useAmiBlockDeviceMappings=false] - ???? Please fill in the details.
**/
local AwsProviderSettings = {
  local this = self,
  useAmiBlockDeviceMappings:: false,

  aws: {
    useAmiBlockDeviceMappings: this.useAmiBlockDeviceMappings,
  },
};

/**
    Creates a new Spinnaker Application object.

    Takes in custom banners and a data-source config. They are parameters in order to conditionally set them if they
    are not null nor empty.

    @example
        local dataSourcesConfig =  application.DataSourcesConfiguration { ... }
        local banners =  [ application.Banner { ... }, ... ]

        ...

        application.NewApplication(banners, dataSourcesConfig) {
            name: 'My Amazing Spinnaker App',
            description: 'Does amazing things with amazing services.',
            email: 'john.doe@shore.com',
        }

    @constructs Application
    @memberof Application
    @name NewApplication

    @param {Array<Banner>} customBanners - An array of Banner objects that will be applied.
    @param {DataSourcesConfiguration} dataSourcesConfig - The DataSourcesConfiguration object to use.

    @return {Application} A new Application.

**/
local NewApplication(customBanners, dataSourcesConfig) = Application {
  [if (customBanners != null && std.length(customBanners) > 0) then 'customBanners']: customBanners,
  // dataSources can't be an empty array, causes "Application Not Found" error in the web app. Otherwise would stick to declarative.
  [if (dataSourcesConfig != null && std.length(dataSourcesConfig) > 0) then 'dataSources']: dataSourcesConfig,
};

/**
    A configuration object for the Spinnaker Application's banner.

    In the Spinnaker UI, this is found in the configuration of the application and labled as "Custom Banners".

    The banner will appear in the Spinnaker Application UI, at the top of the application - near the address bar - on
    all sections of the Spinnaker Application.

    A Spinnaker Application can have multiple Banners.

    @example
        local myBanner = application.Banner {
            text: "This is my custom banner, which will appear in the Spinnaker Application UI."
        }

        ...

        application.Application {
            ...
            customBanners: [myBanner, ...],
            ...
        }
    @property {String} text - The text that the banner will display.
    @property {String} [backgroundColor='var(--color-alert)'] - The background color.
    @property {String} [textColor='var(--color-text-on-dark)'] - The text color.
**/
local Banner = {
  text: error '`text` (String) is  property is required for Banner',

  enabled: true,
  backgroundColor: 'var(--color-alert)',
  textColor: 'var(--color-text-on-dark)',
};

/**
    Map containing the various Spinnaker Application DataSources.

    Each DataSource maps to a section in the Spinnaker UI.

    @example
        application.DataSource.functions

    @enum {String}
    @readonly

    @property {String} executions - Executions DataSource
    @property {String} functions - Functions DataSource
    @property {String} loadBalancers - Load Balancer DataSource
    @property {String} securityGroups - Security Groups DataSource
    @property {String} serverGroups - Server Groups DataSource
**/
local DataSource = {
  executions: 'executions',  // Represented by the "Pipelines" section of the application in the Spinnaker UI.
  functions: 'functions',
  loadBalancers: 'loadBalancers',
  securityGroups: 'securityGroups',
  serverGroups: 'serverGroups',  // Represented by the "Clusters" section of the application in the Spinnaker UI.
};

/**
    An application configuration object holding enabled/disabled DataSources (serverGroups/loadBalancers/executions/...)

    In the Spinnaker UI, this is found in the configuration of the application and labled as "Features".

    Enabling/Disabling DataSources will make their corresponding sections appear/disappear in the Spinnaker UI.

    By default, if not specified, all sections appear in the Spinnaker UI.

    @example
        application.DataSourcesConfiguration {
            enabled: [
                        application.DataSource.serverGroups,
                        application.DataSource.loadBalancers,
                        application.DataSource.executions
                     ],
            disabled: [
                        application.DataSource.securityGroups,
                        application.DataSource.functions
                      ]
        }
    @property {Array<Enum<DataSource>>} enabled - An array of all enabled DataSources string values.
    @property {Array<Enum<DataSource>>} disabled - An array of all disabled DataSources string values.
**/
local DataSourcesConfiguration = {
  enabled: error '`disabled` and `enabled` are required arrays for DataSourcesConfiguration. One of them must have at least one element, the other can be empty.',
  disabled: error '`disabled` and `enabled` are required arrays for DataSourcesConfiguration. One of them must have at least one element, the other can be empty.',
};

// Exposed for public use.
{
  Application:: Application,
  AwsProviderSettings:: AwsProviderSettings,
  Banner:: Banner,
  DataSource:: DataSource,
  DataSourcesConfiguration:: DataSourcesConfiguration,

  NewApplication:: NewApplication,
}
