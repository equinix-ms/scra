package audit

import (
	"go.uber.org/zap/zapcore"
)

func marshalStringArray(array []string) zapcore.ArrayMarshalerFunc {
	return func(enc zapcore.ArrayEncoder) error {
		for _, element := range array {
			enc.AppendString(element)
		}

		return nil
	}
}

func (n *Namespaces) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("networking", n.HostNetworking)
	enc.AddBool("pid", n.HostPID)
	enc.AddBool("ipc", n.HostIPC)
	enc.AddBool("uts", n.HostUTS)
	enc.AddBool("cgroup", n.HostCgroup)

	return nil
}

func (n *NetworkInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("device", n.Device)
	enc.AddArray("addresses", marshalStringArray(n.Addresses))

	return nil
}

func (l *LinuxCapabilities) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddArray("bounding", marshalStringArray(l.Bounding))
	enc.AddArray("effective", marshalStringArray(l.Effective))
	enc.AddArray("inheritable", marshalStringArray(l.Inheritable))
	enc.AddArray("permitted", marshalStringArray(l.Permitted))
	enc.AddArray("ambient", marshalStringArray(l.Ambient))

	return nil
}

func (l *LinuxDevice) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("path", l.Path)
	enc.AddString("type", l.Type)
	enc.AddInt64("major", l.Major)
	enc.AddInt64("minor", l.Minor)

	return nil
}

func (r *Report) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("runtime", r.Runtime)
	enc.AddString("id", r.ID)
	enc.AddString("image", r.Image)
	enc.AddInt("pid", r.PID)
	enc.AddObject("namespaces", &r.Namespaces)
	enc.AddTime("created", r.Created)
	enc.AddArray("mounts", marshalStringArray(r.Mounts))
	enc.AddString("cgroups_path", r.CgroupsPath)
	enc.AddString("status", r.Status)
	enc.AddObject("capabilities", r.Capabilities)

	enc.AddArray("devices", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, device := range r.Devices {
			e.AppendObject(device)
		}
		return nil
	}))

	return nil
}
