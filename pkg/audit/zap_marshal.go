package audit

import (
	"go.uber.org/zap/zapcore"
)

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
	enc.AddArray("addresses", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, address := range n.Addresses {
			e.AppendString(address)
		}
		return nil
	}))

	return nil
}

func (l *LinuxCapabilities) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddArray("bounding", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, capability := range l.Bounding {
			e.AppendString(capability)
		}
		return nil
	}))

	enc.AddArray("effective", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, capability := range l.Effective {
			e.AppendString(capability)
		}
		return nil
	}))

	enc.AddArray("inheritable", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, capability := range l.Inheritable {
			e.AppendString(capability)
		}
		return nil
	}))

	enc.AddArray("permitted", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, capability := range l.Permitted {
			e.AppendString(capability)
		}
		return nil
	}))

	enc.AddArray("ambient", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, capability := range l.Ambient {
			e.AppendString(capability)
		}
		return nil
	}))

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

	enc.AddArray("mounts", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, mount := range r.Mounts {
			e.AppendString(mount)
		}
		return nil
	}))

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
