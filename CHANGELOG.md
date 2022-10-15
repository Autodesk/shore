# Changelog

## [v0.0.9](https://github.com/Autodesk/shore/releases/tag/v0.0.9) (2022-04-05)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.8...v0.0.9)

**Implemented enhancements:**

- Target specific application / Skip stages [\#173](https://github.com/Autodesk/shore/issues/173)
- Upgrade CI/CD to use golang 1.17 [\#161](https://github.com/Autodesk/shore/issues/161)
- JSONNET Renderer - Unit Testing \(planning\) [\#29](https://github.com/Autodesk/shore/issues/29)

**Fixed bugs:**

- Mismatch in passing in parameter blobs from E2E.yml and parameter blobs from exec.yml  [\#175](https://github.com/Autodesk/shore/issues/175)
- Remove Viper as the config manager [\#127](https://github.com/Autodesk/shore/issues/127)
- Shore renders nested variables in an unexpected way [\#122](https://github.com/Autodesk/shore/issues/122)
- Shore render lowercases all parameters, this affects usability. [\#83](https://github.com/Autodesk/shore/issues/83)

**Documentation updates:**

- Release Plan Alpha V1 [\#79](https://github.com/Autodesk/shore/issues/79)

**Closed issues:**

- Migrate to Golang 1.18 [\#181](https://github.com/Autodesk/shore/issues/181)
- add ability to specify the test-suite to be run [\#149](https://github.com/Autodesk/shore/issues/149)
- Sometimes save works, but exec and test-remote get X509 cert errors [\#148](https://github.com/Autodesk/shore/issues/148)
- Update go-jsonnet to `latest` [\#77](https://github.com/Autodesk/shore/issues/77)

## [v0.0.8](https://github.com/Autodesk/shore/releases/tag/v0.0.8) (2021-08-11)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.7...v0.0.8)

**Fixed bugs:**

- Regressions - Any trigger that isn't a pipeline trigger causes a panic [\#145](https://github.com/Autodesk/shore/issues/145)

## [v0.0.7](https://github.com/Autodesk/shore/releases/tag/v0.0.7) (2021-07-15)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.6...v0.0.7)

**Fixed bugs:**

- shore test-remote fails [\#123](https://github.com/Autodesk/shore/issues/123)

## [v0.0.6](https://github.com/Autodesk/shore/releases/tag/v0.0.6) (2021-07-12)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.5...v0.0.6)

**Fixed bugs:**

- exec command Regression in v0.0.5: -p flag is not working and missing exec.yml throws an error [\#137](https://github.com/Autodesk/shore/issues/137)
- error in  project init's pipeline generation in all yamls [\#108](https://github.com/Autodesk/shore/issues/108)

**Documentation updates:**

- code-coverage badge for shore/golang [\#115](https://github.com/Autodesk/shore/issues/115)

**Closed issues:**

- Refactor Init to make use of Go's embedded FS [\#135](https://github.com/Autodesk/shore/issues/135)
- Update `shore project init` with `shore cleanup` files. [\#126](https://github.com/Autodesk/shore/issues/126)

## [v0.0.5](https://github.com/Autodesk/shore/releases/tag/v0.0.5) (2021-06-20)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.4...v0.0.5)

**Implemented enhancements:**

- Shore-CLI - Reverse Pipeline \("Destroy Command"\) [\#34](https://github.com/Autodesk/shore/issues/34)
- Spinnaker Backend - Get pipeline by AppName + PipelineName/PipelineID [\#31](https://github.com/Autodesk/shore/issues/31)

**Documentation updates:**

- New git tag doesn't trigger a build and CI issues [\#99](https://github.com/Autodesk/shore/issues/99)
- Implement Architecture.md file \(Tech-Implementation\) [\#92](https://github.com/Autodesk/shore/issues/92)
- Provide easy development Docker Containers [\#89](https://github.com/Autodesk/shore/issues/89)

**Closed issues:**

- Shore project Integration testing  [\#84](https://github.com/Autodesk/shore/issues/84)
- Distribute a Docker container to run shore without installing the binary. [\#54](https://github.com/Autodesk/shore/issues/54)
- Add missing default files to `shore project init` [\#47](https://github.com/Autodesk/shore/issues/47)

## [v0.0.4](https://github.com/Autodesk/shore/releases/tag/v0.0.4) (2021-04-01)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.3...v0.0.4)

**Fixed bugs:**

- sub-commands payload/render-values flags are handled incoorectly [\#87](https://github.com/Autodesk/shore/issues/87)

**Documentation updates:**

- Provide installation instructions [\#97](https://github.com/Autodesk/shore/issues/97)

## [v0.0.3](https://github.com/Autodesk/shore/releases/tag/v0.0.3) (2021-03-30)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.2...v0.0.3)

**Fixed bugs:**

- Shore Render should render without a `render.yaml` file. [\#81](https://github.com/Autodesk/shore/issues/81)

**Documentation updates:**

- Shore tutorials [\#75](https://github.com/Autodesk/shore/issues/75)
- Create a CONTRIBUTE.md guidelines doc [\#26](https://github.com/Autodesk/shore/issues/26)
- Create a CHANGELOG.MD file [\#25](https://github.com/Autodesk/shore/issues/25)

**Closed issues:**

- Shore-CLI release pipeline setup [\#24](https://github.com/Autodesk/shore/issues/24)

## [v0.0.2](https://github.com/Autodesk/shore/releases/tag/v0.0.2) (2021-03-14)

[Full Changelog](https://github.com/Autodesk/shore/compare/v0.0.1...v0.0.2)

## [v0.0.1](https://github.com/Autodesk/shore/releases/tag/v0.0.1) (2021-03-14)

[Full Changelog](https://github.com/Autodesk/shore/compare/6cf95adbf5e3b939dcf569a3d6cdc0017c3b0f78...v0.0.1)

**Implemented enhancements:**

- Add a polling mechanism for pipeline exec. [\#16](https://github.com/Autodesk/shore/issues/16)

**Fixed bugs:**

- test-remote panics when target pipeline doesn't exist [\#57](https://github.com/Autodesk/shore/issues/57)
- Pass the artifacts property to Spinnaker [\#55](https://github.com/Autodesk/shore/issues/55)
- `shore` produces the wrong error when `main.pipeline.jsonnet` doesn't exist [\#52](https://github.com/Autodesk/shore/issues/52)
- test-remote bug: produces panic  [\#41](https://github.com/Autodesk/shore/issues/41)
- Is it common to pass logger's around? [\#30](https://github.com/Autodesk/shore/issues/30)
- Nested bug fix: add application injection for child pipelines [\#23](https://github.com/Autodesk/shore/issues/23)

**Documentation updates:**

- Automated API Docs Generation [\#61](https://github.com/Autodesk/shore/issues/61)

**Closed issues:**

- another test issue for changelog test [\#76](https://github.com/Autodesk/shore/issues/76)
- Implement a CLI interface to pass parameters to `exec`. [\#69](https://github.com/Autodesk/shore/issues/69)
- test issue  [\#66](https://github.com/Autodesk/shore/issues/66)
- Jsonnet-bundler review [\#64](https://github.com/Autodesk/shore/issues/64)
- Delete `exec --save` command [\#56](https://github.com/Autodesk/shore/issues/56)
- Copy paste Skipper lib\_adfs as a new github project \(shared-library\) [\#53](https://github.com/Autodesk/shore/issues/53)
- Armory-Specific-shared-library [\#50](https://github.com/Autodesk/shore/issues/50)
- ADSK-Spinnaker-shared-library [\#49](https://github.com/Autodesk/shore/issues/49)
- Spinnaker-stdlib [\#48](https://github.com/Autodesk/shore/issues/48)
- Implement `shore get` wrapper for downloading 3rd party libraries [\#45](https://github.com/Autodesk/shore/issues/45)
- Add logging everywhere! [\#15](https://github.com/Autodesk/shore/issues/15)
- Create a basic Intergation testing suite. [\#14](https://github.com/Autodesk/shore/issues/14)
- Allow customers to pass execution parameters to the Exec command [\#13](https://github.com/Autodesk/shore/issues/13)
- Allow customers to pass arguments when rendering a pipeline [\#8](https://github.com/Autodesk/shore/issues/8)
- Test out how multiple versions of the same package are handled in multiple depedant projects. [\#4](https://github.com/Autodesk/shore/issues/4)
- Implement the nested pipeline save on the FW level [\#3](https://github.com/Autodesk/shore/issues/3)
- Implement a basic nested pipeline example. [\#2](https://github.com/Autodesk/shore/issues/2)
- Implement a basic managed infra example \(with either wait stages or very basic TF stages\) [\#1](https://github.com/Autodesk/shore/issues/1)



