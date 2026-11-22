#!/bin/sh
# elmos demo: load every module in the shared folder, with a result line each.
echo '  [GUEST] Processing module queues...'
for ko in $(find /mnt/modules -name '*.ko' -type f 2>/dev/null); do
    name=$(basename "$ko" .ko)
    if insmod "$ko"; then
        echo "  [GUEST] $name loaded OK"
    else
        echo "  [GUEST] $name FAILED"
    fi
done
if [ -r /sys/kernel/debug/kdal/demand ] || { mount -t debugfs debugfs /sys/kernel/debug 2>/dev/null && [ -r /sys/kernel/debug/kdal/demand ]; }; then
    echo '  [GUEST] KDAL resource demand (/sys/kernel/debug/kdal/demand):'
    cat /sys/kernel/debug/kdal/demand
fi
