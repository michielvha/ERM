# Edge Resource Manager

This repository details my journey on the creation of a resource management API using Golang. I based the name on the ``ARM/AzureResourceManager``. I'll be using it in my cloud platform based on edge/IoT devices.  

[//]: # (But why reinvent the wheel and create the API from scratch? Because I want to learn how you make one and most importantly how you would scale such an endeavour. If Edge Cloud ever becomes a thing I'll hire some russians to do it for me until then I need to grasp all the concepts and layers involved.)

I've been playing with the idea of a distributed cloud system for a while. Such an API might help in building such a platform.

## Next steps


- Provide ability to create new users => Create admin panel with proper RBAC, that allows you to configure new users.
- Create groups for RBAC
- Have a protect endpoint execute armbian build framework based on input headers.


### Enhancements
- **Integrate with a real OAuth2 provider** like Auth0 or Keycloak for better security and compliance.
- **Make it run code** I want to use this API to be in front my infrastructure. The first functionality it should do is create custom images with the ARM build framework.
- **Enhance JWT claims** to include user roles, scopes, or additional data.
- **Implement token refresh logic** if necessary.


### Change log

**0.1.30**
- Binary now takes env vars to setup connection with DB, if none are provided it will use the defaults set in code. This is just a fallback for dev env and should never be used in production.

**0.1.29**
- Added a default admin user to database. DB is migrated with placeholder password, after the initial migrate a function is called to update the password to a secure one based on env var.

**0.1.28**
- modified login function to actually care about users, not just the example workflow we had

**0.1.27**
- Created a migration to provision a user schema.

**0.1.26**
- Added a database to store users & credentials.
- Created migrations to provision a user schema. & embedded them into the binary using iofs.