# Redis Package

## Mappings

The following Redis mappings are maintained in the redis instance.

| Prefix   | Key                 | Value                                            |
| :---     |    :----:           |         :---                                    |
| usr2vgr- | Vault User Alias    | Set of Vault Group Names                         |
| vgr2pol- | Vault Group Name    | Set of Policy Names                              |
| usr2ent- | Vault User Alias    | Entity ID                                        |
| ent2pol- | Entity ID           | Set of Policy Names                              |
| usr2egr- | Vault User Alias    | Set of External Group Names                      |
| egr2rol- | External Group Name | Set of Vault Roles                               |
| rol2pol- | Role                | Set of Policy Names                              | 
| pat2pol- | Vault Path          | Set of Encodings of Policy Name and capabilities |
| pol2vgr- | Vault Policy Name   | Set of Vault Group Names                         |
| vgr2usr- | Vault Group Name    | Set of Vault User Aliases                        |
| pol2ent- | Policy Name         | Set of Entity IDs                                | 
| ent2usr- | Vault User Alias    | Entity ID                                        |    
| pol2rol- | Policy Name         | Set of Roles                                     |
| rol2egr- | Role                | Set of External Group Names                      |
| egr2usr- | External Group Name | Set of Vault User Aliases                        |
| pol2sec- | Policy Name         | Set of Vault Paths                               |    
