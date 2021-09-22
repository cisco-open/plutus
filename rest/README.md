# REST API

## Specification

Prefix all the paths with /api/v1 to use the first version.

1. **Namespaces**

    Fetches all the namespaces aailable to plutus

   * **URL**

        /namespaces

   * **Method:**

        `GET`

   * **Success Response:**  
     **Code:** 200  
     **Body**  

     ```golang
      type NamspacesResponse struct {
        Namespaces []string `json:"namespaces"`
      } 
     ``` 

       <details>
        <summary> <b>Sample</b> </summary>

        ```json
        { 
          "namespaces" : ["ns1", "ns2"] 
        }
        ```
       </details>
<br>       

2. **Last Refreshed**

    Fetches the time that has passed since the redis instance was refreshed

   * **URL**
  
        /lastRefresh

   * **Method:**

        `GET`

   * **Success Response:**  
     **Code:** 200  
     **Type**  

     ```golang
      type LastRefreshResponse struct {
    	  PassedTime int `json:"passedTime"`
      } 
     ``` 

      <details>
      <summary> <b>Sample</b> </summary>

        ```json
        { 
          "passedTime" : 15
        }
        ```

      </details>
<br>

3. **UI Handler**

    Redirects to the UI with the baseURL set

   * **URL**

     /ui

   * **Method:**

     `GET`

   * **Success Response:**  
     **Code:** 302

4. **Vault Path**

    Returns a JSON object containg the vault aliases that can access the given path.

   * **URL**

     /path

   * **Payload**

    ```json  
     {
       "namespace": "ns",
       "path": "some/vault/path" 
     }
    ```

   * **Method:**

     `POST`

   * **Success Response:**  
     **Code:** 200  
     **Type**  

     ```golang
      type Response struct {
        Users []Users `json:"users"`
      }

      type Policies struct {
        Name         string   `json:"name"`
        Roles        []string `json:"roles"`
        Groups       []string `json:"groups"`
        Entities     []string `json:"entities"`
        Capabilities []string `json:"capabilities"`
      }

      type Users struct {
        Username     string     `json:"username"`
        Capabilities []string   `json:"capabilities"`
        Policies     []Policies `json:"policies"`
      }
     ```
     <details>
     <summary> <b>Sample</b> </summary>
     
        ```json
        {
          "users": [
            {
              "username":"pranbans",
              "capabilities": ["read", "list"],
              "policies": [
                {
                  "name" : "auth-readonly",
                  "roles": ["admin", "devops"],
                  "groups": ["developers"],
                  "entities": ["000ae7e3-6ed4-6693-4527-705ecafbbc42"],
                  "capabilities": ["read", "list"]
                }
              ] 
            }
          ]
        }
        ```

     </details>

<br>

5. **Vault Username(Alias)**

    Returns a JSON object containg the information about the users vault presence.

   * **URL**

     /user

   * **Payload**

    ```json  
     {
       "namespace": "ns",
       "username": "someVaultUsername" 
     }
    ```

   * **Method:**

     `POST`

   * **Success Response:**  
     **Code:** 200  
     **Type**  

     ```golang
      type Response struct {
        Groups   []string   `json:"groups"`
        Roles    []string   `json:"roles"`
        Policies []Policies `json:"policies"`
        Paths    []Paths    `json:"paths"`
      }
      
      type Policies struct {
        Name     string   `json:"name"`
        Roles    []string `json:"roles"`
        Groups   []string `json:"groups"`
        Entities []string `json:"entities"`
      }

      type Paths struct {
        Path         string   `json:"path"`
        PolicyName   string   `json:"policy_name"`
        Capabilities []string `json:"capabilities"`
      }
     ```
     <details>
     <summary> <b>Sample</b> </summary>

        ```json
        {
          "groups": ["g1", "g2"],
          "roles": ["admin", "devops"],
          "policies": [
            {
              "name" : "auth-readonly",
              "roles": ["admin", "devops"],
              "groups": ["developers"],
              "entities": ["000ae7e3-6ed4-6693-4527-705ecafbbc42"]
            }
          ],
          "paths": [
            {
              "path": "som/vault/path",
              "policy_name": "group-readonly",
              "capabilities": ["read", "list"]
            }
          ] 
        }
        ```
    </details>