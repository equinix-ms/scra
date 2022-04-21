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
	return enc.AddArray("addresses", marshalStringArray(n.Addresses))
}

func (l *LinuxCapabilities) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := enc.AddArray("bounding", marshalStringArray(l.Bounding))
	if err != nil {
		return err
	}

	err = enc.AddArray("effective", marshalStringArray(l.Effective))
	if err != nil {
		return err
	}

	err = enc.AddArray("inheritable", marshalStringArray(l.Inheritable))
	if err != nil {
		return err
	}

	err = enc.AddArray("permitted", marshalStringArray(l.Permitted))
	if err != nil {
		return err
	}

	return enc.AddArray("ambient", marshalStringArray(l.Ambient))
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
	enc.AddString("namespace", r.Namespace)
	enc.AddString("image", r.Image)
	enc.AddInt("pid", r.PID)
	err := enc.AddObject("host_namespaces", &r.HostNamespaces)
	if err != nil {
		return err
	}

	enc.AddTime("created", r.Created)
	err = enc.AddArray("mounts", marshalStringArray(r.Mounts))
	if err != nil {
		return err
	}

	enc.AddString("cgroups_path", r.CgroupsPath)
	enc.AddString("status", r.Status)
	err = enc.AddObject("capabilities", r.Capabilities)
	if err != nil {
		return err
	}

	return enc.AddArray("devices", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, device := range r.Devices {
			err := e.AppendObject(&device)
			if err != nil {
				return err
			}
		}
		return nil
	}))
}
