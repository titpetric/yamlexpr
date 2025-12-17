### Basic String Interpolation

Interpolate variable values into strings using `${}` syntax.

**Input:**

```yaml
name: "World"
greeting: "Hello ${name}"
message: "Welcome back, ${name}!"
```

**Output:**

```yaml
name: "World"
greeting: "Hello World"
message: "Welcome back, World!"
```

### Nested Path Interpolation

Access nested object properties with dot notation in interpolation.

**Input:**

```yaml
user:
  name: "alice"
  email: "alice@example.com"
  profile:
    title: "Engineer"
messages:
  welcome: "Welcome ${user.name}"
  contact: "Reach ${user.name} at ${user.email}"
  about: "${user.name} is a ${user.profile.title}"
```

**Output:**

```yaml
user:
  name: "alice"
  email: "alice@example.com"
  profile:
    title: "Engineer"
messages:
  welcome: "Welcome alice"
  contact: "Reach alice at alice@example.com"
  about: "alice is a Engineer"
```

### Expression Interpolation

Evaluate expressions and calculations within `${}` syntax.

**Input:**

```yaml
base_price: 100
discount_percent: 15
quantity: 5
pricing:
  unit_price: "${base_price}"
  discount_amount: "${base_price * discount_percent / 100}"
  total: "${(base_price * (100 - discount_percent) / 100) * quantity}"
  summary: "Total for ${quantity} units: ${(base_price * (100 - discount_percent) / 100) * quantity}"
```

**Output:**

```yaml
base_price: 100
discount_percent: 15
quantity: 5
pricing:
  unit_price: "100"
  discount_amount: "15"
  total: "425"
  summary: "Total for 5 units: 425"
```

### Multiple Interpolations

Multiple variables and expressions in a single string.

**Input:**

```yaml
first_name: "John"
last_name: "Doe"
version: "1.0.0"
build_number: 42
release_date: "2024-01-15"
url_base: "https://example.com"
artifact: "app-${version}-build${build_number}.tar.gz"
full_name: "${first_name} ${last_name}"
download_url: "${url_base}/releases/${version}/${artifact}"
release_info: "${full_name} released ${artifact} on ${release_date}"
```

**Output:**

```yaml
first_name: "John"
last_name: "Doe"
version: "1.0.0"
build_number: 42
release_date: "2024-01-15"
url_base: "https://example.com"
artifact: "app-1.0.0-build42.tar.gz"
full_name: "John Doe"
download_url: "https://example.com/releases/1.0.0/app-1.0.0-build42.tar.gz"
release_info: "John Doe released app-1.0.0-build42.tar.gz on 2024-01-15"
```

### Interpolation in For Loops

Interpolate loop variables and context values within for loop templates.

**Input:**

```yaml
namespace: "production"
replicas: 3
services:
  - for: service in ["api", "worker", "cache"]
    name: "${service}"
    image: "myrepo/${service}:latest"
    replicas: ${replicas}
    full_name: "${namespace}-${service}"
```

**Output:**

```yaml
namespace: "production"
replicas: 3
services:
  - name: "api"
    image: "myrepo/api:latest"
    replicas: 3
    full_name: "production-api"
  - name: "worker"
    image: "myrepo/worker:latest"
    replicas: 3
    full_name: "production-worker"
  - name: "cache"
    image: "myrepo/cache:latest"
    replicas: 3
    full_name: "production-cache"
```
