virsh pool-define-as storage_pool --type dir --target $(pwd)/storage_pool
virsh pool-start --build storage_pool