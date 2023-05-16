MERGE(os:OrderStatus {uid: 'c5ef9d00-ac38-44c1-b48a-fde0d7095c56'})
SET os.created = timestamp(), os.updated = timestamp(), os.code = 'ORDER_COMPLETED', os.name = 'Order completed', os.sortOrder = 5;