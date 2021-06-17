# DefaultApi

All URIs are relative to *http://localhost:8000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**messageList**](DefaultApi.md#messageList) | **GET** /messages | 


<a name="messageList"></a>
# **messageList**
> SelectedMessage messageList(pageSize, pageStartIndex, fields)



    List all or a subset of existing messages.

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pageSize** | **Integer**| Limits the number of returned rows. | [optional] [default to null]
 **pageStartIndex** | **Integer**| Determines query page number of a given size pageSize. | [optional] [default to null]
 **fields** | **String**| Limits the returned fields to those specified here. | [optional] [default to null]

### Return type

[**SelectedMessage**](../Models/SelectedMessage.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

