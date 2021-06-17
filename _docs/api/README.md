# Documentation for Title

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost:8000*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*DefaultApi* | [**messageList**](Apis/DefaultApi.md#messagelist) | **GET** /messages | List all or a subset of existing messages.
*MessageApi* | [**messageCreate**](Apis/MessageApi.md#messagecreate) | **POST** /messages | Create a new message.
*MessageApi* | [**messageDeleteById**](Apis/MessageApi.md#messagedeletebyid) | **DELETE** /messages/{id} | Delete a message.
*MessageApi* | [**messageGetById**](Apis/MessageApi.md#messagegetbyid) | **GET** /messages/{id} | Retrieve an existing message.
*MessageApi* | [**messageUpdateById**](Apis/MessageApi.md#messageupdatebyid) | **PUT** /messages/{id} | Update a message.


<a name="documentation-for-models"></a>
## Documentation for Models

 - [ErrorResponse](./Models/ErrorResponse.md)
 - [Message](./Models/Message.md)
 - [MessageModify](./Models/MessageModify.md)
 - [SelectedMessage](./Models/SelectedMessage.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

All endpoints do not require authorization.
