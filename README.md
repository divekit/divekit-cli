# divekit-cli

**This is work in progress, but the Patch Tool (`divekit patch` command) is robust and usable.** 
No other subcommands are usable as of now. 

Learn more about the CLI and other divekit tools in the [Divekit documentation](https://divekit.github.io/docs/cli/)
... but the documentation there is outdated with regard to the CLI. So this README acts as a temporary documentation, 
to be migrated later.


## Assumptions for using the CLI 

The CLI assumes the following aspects to work properly: 
- The following repos need to be cloned under **their original names**, and in the same parent "git" directory:
  - [divekit-cli](https://github.com/divekit/divekit-cli)
  - [divekit-automated-repo-setup](https://github.com/divekit/divekit-automated-repo-setup)
  - [divekit-repo-editor](https://github.com/divekit/divekit-repo-editor)
- You need to have `npm` installed on your machine

That's it. Neither a specific IDE nor Go need to be installed.


## What the Patch Tool does

The advantage of the patch tool is that you omit all the error-prone manual copy-pasting between two tools
(ARS and Repo Editor). The CLI command `divekit patch` automatically performs the following steps:
1. Read all necessary settings from the origin repo (see glossary below)
2. Define the config settings for ARS, based on input flags / parameters and the origin repo information
3. Run the [ARS](#ars) "locally" to produce the individualized patch files in the local file system 
4. Copy the localized files in the appropriate input directories in the [Repo Editor](#repo-editor)
5. Configure the Repo Editor settings based on the input flags / parameters and the origin repo information
6. Run the Repo Editor




## How to call the Patch Command via CLI

Let's assume that you want to patch buggy unit test file called `E2WhateverTests.java`, and the `pom.xml`
where a dependency is missing. Your origin repo is called `st2-m3-origin`.
Let's further assume that you have two [distributions](#distribution):
 
- The standard "milestone" distribution, which contains all the student campus ids and UUIDs, and
  the Gitlab groups for the students. 
- A "test" distribution with just two repos, using supervisor campus ids.

As a first step, open a shell and go to your `divekit-cli` repo dir. (You can run the command from any other location, 
just make sure you have divekit.exe in your path, e.g. by copying it somewhere.)

It is recommended to do a trial run with the "test" distribution, where you patch only the two
supervisors' repos. This is done by the following command:
```
divekit patch -m <my-local-git-dir> -o st2-m3-origin -d test E2WhateverTests.java pom.xml`
```
Note that you don't have to specify the precise path of your files. The patch tool searches relevant locations
in the origin repo to find the full path, and gives an error if there are multiple files matching that name. You
should check in the supervisors' repos if the patch delivered the expected results. 

If this was successful, you can patch the student repos:
```
divekit patch -m <my-local-git-dir> -o st2-m3-origin E2WhateverTests.java pom.xml`
```
The only difference is that the `-d test` option is missing - the tool assumes that the student distribution
is called `milestone`. If it has a different name, just use `-d` accordingly.


## Documentation for flags and parameters

The best way is to call `divekit patch -h`, then you get a brief documentation of available flags.


## Glossary

#### Origin (repo)

The repo where the milestone exercise is defined in. We assume a certain structure. Essential is
the folder `.divekit` on top level, where all information are stored in.

#### Distribution

The list of users for which individualized repos are created. Default distributions are:
- milestone: with the actual students, creating repos in the `student` group
- test: using the supervisor(s) campusIds, creating repos in the `staff` group

But there may be more distributions. Distributions are stored in `.divekit_norepo\distributions` in the origin
repo, like this:

```
.divekit_norepo
    distributions
        milestone
            individual_repositories_03-04-2023 01-39-08.json
            repositoryConfig.json
        test
            individual_repositories_test_ecom.json
            repositoryConfig.json
```

#### ARS

Abbreviation for `divekit-automated-repo-setup`, the core Divekit tool that produces individualized 
files for student exercises, practicals, or exams. See 
[divekit-automated-repo-setup](https://github.com/divekit/divekit-automated-repo-setup).


#### Repo Editor

Despite the name, this is actually a patch tool, replacing certain (individualized) files in the students'
milestone repos. See [divekit-repo-editor](https://github.com/divekit/divekit-repo-editor).



## INTERNAL - Development of the "patch" subcommand - how to test

(this is just temporary during development, will be removed later)

As test data, I have used the following origin repo:
- st2-m0-origin (https://git.st.archi-lab.io/staff/st2/ss22/m0/st2-m0-origin)
  - In the .divekit_norepo, the necessary files have been added
  - There are two test files, one open and one hidden (test repo only). I have added a useless (but not
    confusing) individualized comment on top of the first test in each test file.