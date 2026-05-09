-- Seed default permissions
INSERT INTO permissions (id, name, description) VALUES
    (gen_random_uuid(), 'products:read',       'View products'),
    (gen_random_uuid(), 'products:create',     'Create products'),
    (gen_random_uuid(), 'products:update',     'Update products'),
    (gen_random_uuid(), 'products:delete',     'Delete products'),
    (gen_random_uuid(), 'categories:read',     'View categories'),
    (gen_random_uuid(), 'categories:create',   'Create categories'),
    (gen_random_uuid(), 'categories:update',   'Update categories'),
    (gen_random_uuid(), 'categories:delete',   'Delete categories'),
    (gen_random_uuid(), 'orders:read',         'View orders'),
    (gen_random_uuid(), 'orders:update',       'Update order status'),
    (gen_random_uuid(), 'inventory:read',      'View inventory'),
    (gen_random_uuid(), 'inventory:update',    'Adjust inventory'),
    (gen_random_uuid(), 'customers:read',      'View customer accounts'),
    (gen_random_uuid(), 'admins:read',         'View admin accounts'),
    (gen_random_uuid(), 'admins:create',       'Create admin accounts'),
    (gen_random_uuid(), 'admins:update',       'Update admin accounts'),
    (gen_random_uuid(), 'roles:manage',        'Manage roles and permissions'),
    (gen_random_uuid(), 'dashboard:read',      'View dashboard analytics');

-- Seed default roles
INSERT INTO roles (id, name, description) VALUES
    ('00000000-0000-0000-0000-000000000001', 'super_admin',       'Full access to everything'),
    ('00000000-0000-0000-0000-000000000002', 'manager',           'Manage products, orders, and inventory'),
    ('00000000-0000-0000-0000-000000000003', 'inventory_staff',   'Manage inventory only'),
    ('00000000-0000-0000-0000-000000000004', 'support',           'View orders and customers');

-- super_admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions;

-- manager permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000002', id FROM permissions
WHERE name IN (
    'products:read', 'products:create', 'products:update',
    'categories:read', 'categories:create', 'categories:update',
    'orders:read', 'orders:update',
    'inventory:read', 'inventory:update',
    'customers:read', 'dashboard:read'
);

-- inventory_staff permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000003', id FROM permissions
WHERE name IN ('inventory:read', 'inventory:update', 'products:read');

-- support permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000004', id FROM permissions
WHERE name IN ('orders:read', 'customers:read', 'products:read');
