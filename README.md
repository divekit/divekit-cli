# divekit-cli

Learn more about the CLI and other divekit tools in
the [Divekit documentation](https://divekit.github.io/docs/cli/).

This is a temporary documentation, to be added somewhere else later.

## Glossary:


### Origin (repo) 

The repo where the milestone exercise is defined in. We assume a certain structure. Essential is
the folder `.divekit` on top level, where all information are stored in.

### Distribution

The list of users for which individualized repos are created. Default distributions are:
- milestone: with the actual students, creating repos in the `student` group
- test: using the supervisor(s) campusIds, creating repos in the `staff` group

But there may be more distributions. Distributions are stored in `.divekit\distributions` with the filename
`<name>.repositoryConfig.json`. Example: 

```
.divekit
    distributions
        milestone.repositoryConfig.json
        test.repositoryConfig.json
```


## Global Flags

see divekit.go

## Patch command

As test data, I have used the following origin repo: 
- st2-m0-origin (https://git.st.archi-lab.io/staff/st2/ss22/m0/st2-m0-origin)
    - In the .divekit_norepo, the necessary files have been added
    - There are two test files, one open and one hidden (test repo only). I have added a useless (but not 
      confusing) individualized comment on top of the first test in each test file.
   

