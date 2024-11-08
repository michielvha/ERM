# Edge Resource Manager

This repository details my journey on the creation of a resource management API using Golang. I based the name on the ``ARM/AzureResourceManager``. I'll be using it in my cloud platform based on edge/IoT devices.  

[//]: # (But why reinvent the wheel and create the API from scratch? Because I want to learn how you make one and most importantly how you would scale such an endeavour. If Edge Cloud ever becomes a thing I'll hire some russians to do it for me until then I need to grasp all the concepts and layers involved.)

I've been playing with the idea of a distributed cloud system for a while. Such an API might help in building such a platform.

### Enhancements
- **Integrate with a real OAuth2 provider** like Auth0 or Keycloak for better security and compliance.
- **Add a database** Add a database to store users & credentials.
- **Make it run code** I want to use this API to be in front my infrastructure. The first functionality it should do is create custom images with the ARM build framework.
- **Enhance JWT claims** to include user roles, scopes, or additional data.
- **Implement token refresh logic** if necessary.

