package main

import (
	"fmt"

	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
  wapc "github.com/wapc/wapc-guest-tinygo"
)

func validate(payload []byte) ([]byte, error) {
	settings, err := NewSettingsFromValidationReq(payload)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	svc_type := gjson.GetBytes(
		payload,
		"request.object.spec.type")
	if !svc_type.Exists() {
		logger.Warn("cannot read svc type: accepting request")
		return kubewarden.AcceptRequest()
	}

  svc_ns := gjson.GetBytes(
    payload,
    "request.object.metadata.namespace")
	if !svc_ns.Exists() {
		logger.Warn("cannot read svc NS: accepting request")
		return kubewarden.AcceptRequest()
	}

  existing_ns_list, err := wapc.HostCall("kubernetes", "namespaces", "list", []byte(""))
  if err != nil {
    return kubewarden.RejectRequest(
      kubewarden.Message(err.Error()),
      kubewarden.Code(400))
  }
  logger.Info("Got existing NS List: ")
  logger.Info(string(existing_ns_list))
  existing_ns := gjson.GetBytes(
    existing_ns_list,
    "items")

  // res will hold pointer to namespace
  var res *gjson.Result
  existing_ns.ForEach(func(key, value gjson.Result) bool {
    ns_name := value.Get("metadata.name")
    if ns_name.String() == svc_ns.String() {
      logger.Info("FOUND NS " + ns_name.String())
      res = &value
    }
    return true
  })
  if res == nil {
    return kubewarden.RejectRequest(
      kubewarden.Message(fmt.Sprintf("NS not found")),
        kubewarden.NoCode)
  }

 // go through annotations of NS and check any of them is in settings
  set_annos := res.Get("metadata.annotations")
  logger.Info("Annos of existing NS "+set_annos.String())
  if !set_annos.Exists(){
    return kubewarden.RejectRequest(
      kubewarden.Message(fmt.Sprintf("NS has no annotations")),
        kubewarden.NoCode)
  }
  found := false
  set_annos.ForEach(func(key, value gjson.Result) bool {
    logger.Info("if "+settings.AllowNodePortAnnotations.String()+" contains "+key.String())
    if settings.AllowNodePortAnnotations.Contains(key.String()){
      found = true
    }
    return true
  })
  if found {
    return kubewarden.AcceptRequest()
  }

  return kubewarden.RejectRequest(
    kubewarden.Message(fmt.Sprintf("Missing required annotations on NS : ", settings.AllowNodePortAnnotations)),
      kubewarden.NoCode)
}
