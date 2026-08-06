Statement: The company X now is having 10k servers. Let's build a system to manage the on/off statuses of all of these servers. The broad idea is like this: Instead of normally manually check the servers, users will send API requests to your server management system (Lte's call it SMS-X for abbreviation), SMS-X somehow will serve information on servers statuses

Requirements:

A. Functional requirements:

1. A server has some fundamental infor:

- server_id: Unique ID for each server
- server_name: Name of servers. It's unique too
- status: On/off
- created_time: When the server was created
- last_updated: Last time server got updated
- ipv4: IP address of the server
- Any extra info u needs to, to help with functionalities, you can define more yourself

2. Functionalities to the system:

- Status check of the servers: SMS-X status check servers in the management list and update statuses to the centralized system (They dont say which, but I believe SMS-X in this case)
- Manage server:
    - Create server: Allow user to create a server will full info. 
        - Input: Server info. 
        - Output: Server creation status
    
    - View server: Get list of servers. Can have more filters if user add filters. The list must be paginated, and allow sorting by field(s)
        - Input: 
            - Filter(optional)
            - From/to: Pagination info
            - Sort: Sort field, and either ascending/descending order
        - Output:
            - Number of suitable servers
            - List of suitable servers
    - Update server: Update server information. You can update anything, except server_id (There's also server name, but they don't said that I can't quite update server's name, so I guess it's updatable, but you have to have a mechanism to check, to allow updating name from X to Y in the moment the user send request, or nah. That's smt to think)
        - Input:
            - server_id: ID of the server to be updated
            - update_data: Information that you want to use to update the server
        - Output:
            - Information of the server post update
    - Delete server: Delete a server from servers management list
        - Input: server_id for deletion
        - Output: Deletion result
    - Import servers: Allow create a list of servers from an excel file (They don't say what here, but my guess is that, from .csv file. Just assumptions). Skip server IDs or server names that already exist in the system
        - Input: File to import from
        - Output:
            - Number of successfully imported servers
            - List of successfully imported servers
            - Number of failed imported servers
            - List of failed imported servers
    - Export servers: Allow exporting servers list to an excel file as output
        - Input: 
            - Filter (optional)
            - From/to (pagination)
            - Sort field and which order (ascend/descending)
        - Output: Excel file with suitable servers, as the result of the above query
    - Reporting: You have to make sure the system sends email to admin (Could be more than one idk ??) to report on server status the prior day (So let's say, today is 20th May 2026, right ? All the admin(s) should have received an email on system statuses in the morning to start off). The email should have:
        - Number of servers in the system
        - Number of servers that were on
        - Number of servers that were off
        - Average uptime of all servers
        - Aside from those above info, you need to also build an API yourself, so that, let's say, an admin wants to have that report himself, could call API and get the report and email sends to his admin account or smt. It looks like this:
            - Input: 
                - The beginning and end date of report (So like, a range day you know. It's not just one day, it's not just yesterday, it could be aggregation of multiple days in the past)
                - Admin email(s) (Could be more than one email ?? Project desc didn't say anything, so I guess we need to make our own thoughts here)
            - Output: Send email to the admin(s)

B. Non-functional requirements:

- Coding convention:
    - Must use OpenAPI (So yeah, writing specs and stuff)
    - Must have unit tests with code coverage >= 90%
    - Use modules/libraries to interact with Database and prevent all possible SQL injections
    - Output of APIs must defines clearly all cases of successes/errors with appropriate error codes and descriptions
    - Output logs to files and you must have logrotate (Not sure what this means, maybe it's some log on Linux stuff)
- Authen/Authorization:
    - All APIs must be authenticated/authorized by JWT
    - Each API must have their own scope
- Technology:
    - Use Redis cache for performance optimizations when necessary
    - Use Elasticsearch to calculate uptime of a server
    - Use Postgres for DB
    - Some other techs that you feel if needed, and it's not limited