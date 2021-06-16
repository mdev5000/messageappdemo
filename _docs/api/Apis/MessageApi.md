# MessageApi

All URIs are relative to *http://localhost:8000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**messageCreate**](MessageApi.md#messageCreate) | **POST** /messages | 
[**messageDeleteById**](MessageApi.md#messageDeleteById) | **DELETE** /messages/{id} | 
[**messageGetById**](MessageApi.md#messageGetById) | **GET** /messages/{id} | 
[**messageUpdateById**](MessageApi.md#messageUpdateById) | **PUT** /messages/{id} | 


<a name="messageCreate"></a>
# **messageCreate**
> messageCreate(MessageModify)



    Create a new message.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **MessageModify** | [**MessageModify**](../Models/MessageModify.md)|  | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="messageDeleteById"></a>
# **messageDeleteById**
> messageDeleteById(id)



    Delete a message.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Long**| Message Id | [default to null]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

<a name="messageGetById"></a>
# **messageGetById**
> Message messageGetById(id)



    Retrieve an existing message.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Long**| Message Id | [default to null]

### Return type

[**Message**](../Models/Message.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="messageUpdateById"></a>
# **messageUpdateById**
> messageUpdateById(id, MessageModify)



    Update a message.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Long**| Message Id | [default to null]
 **MessageModify** | [**MessageModify**](../Models/MessageModify.md)|  | [optional]

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

