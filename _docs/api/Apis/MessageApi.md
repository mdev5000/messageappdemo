# MessageApi

All URIs are relative to *http://localhost:8000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**messagesIdGet**](MessageApi.md#messagesIdGet) | **GET** /messages/{id} | 
[**messagesPost**](MessageApi.md#messagesPost) | **POST** /messages | 


<a name="messagesIdGet"></a>
# **messagesIdGet**
> Message messagesIdGet(id)



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

<a name="messagesPost"></a>
# **messagesPost**
> messagesPost()



    Create a new message.

### Parameters
This endpoint does not need any parameter.

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

