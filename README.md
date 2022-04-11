# SCRA: Static container runtime audit

A tool for static analysis of the contents of your container runtime. This tool
helps you identify containers with "wild" privileges, strange mounts, host 
namespace access, etc.

It does not monitor the dynamic behaviour of containers, only their definition.

scra is licensed under GPLv3