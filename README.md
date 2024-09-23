# coredns-libvirt

This plugin is based on the [coredns-libvirt](https://github.com/michaelbeaumont/coredns-libvirt) implementation.
But this implementation is working with the libvirtd and not with the dhcp leases files.

## Usage

Currently this plugin can only be used in the `guest` mode.

### Guest name

The functionality of `libvirt guest` is analogous to the `libvirt_guest` `nss`
plugin, where we look for a match on the name of the libvirt domain, not
necessarily a hostname.

### Connection Uri

The connection uri to connect with libvirtd. If not set the default uri `qemu:///system`.

### Network

The network name where the domains are connected to. If not set the default network name is `default`.

### Name map

IF you want to map a domain name to a diffrent dns name you can use this parameter.
e.g. `name_map nextcloud cloud` to response to dns request for `cloud` with the ip address from domain `nextcloud`

### Zones

If your zone isn't root `.`, you'll likely want to include the `trim_suffix`
directive so you search for the correct name in your guests.

### Filtering by network

If only some of the IPs assigned to the guests are reachable, you can filter
them with the `keep` directive.

## Example

```
vm.network:1053 {
  libvirt guest {
    connect_uri qemu:///system
    network default
    trim_suffix vm.network
    keep 10.101.0.0/24
    name_map nextcloud cloud
    name_map jenkins gitlab_runner
  }
}
```
