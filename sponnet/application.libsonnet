{
  newApplication:: {
    cloudProviders: 'aws, kubernetes',
    // Empty dataSources causing "Application Not Found" error in the web app? :thinking:
    //dataSources: {
    //  disable: [],
    //  enabled: [],
    //},
    description: error 'description required',
    email: error 'email required',
    name: error 'name required',
    providerSettings: {
      aws: { useAmiBlockDeviceMappings: false },
    },
  },
}
