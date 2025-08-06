# 🚀 ABDD - API Business Driven Development

> **Declarative API testing that speaks your business language**

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/davesavic/abdd)](https://goreportcard.com/report/github.com/davesavic/abdd)

ABDD is a powerful API testing tool that bridges the gap between business requirements and technical implementation. Write your API tests in human-readable YAML files that describe **what** your API should do, not **how** to test it.

## ✨ What Makes ABDD Special?

- **🎯 Business-First**: Tests read like business requirements, not technical specifications
- **🔗 Smart Dependencies**: Automatically handles test execution order based on dependencies
- **🎲 Rich Fake Data**: Built-in integration with 310+ fake data generators
- **🔄 Variable Extraction**: Extract values from responses to use in subsequent tests
- **⚡ Zero Configuration**: Get started with a single command
- **📝 YAML-Based**: Familiar, readable syntax that non-developers can understand
- **🛠️ Shell Integration**: Execute custom commands as part of your test flow

## 📋 Table of Contents

- [What is Business Driven Development?](#what-is-business-driven-development)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Core Concepts](#core-concepts)
- [Writing Tests](#writing-tests)
  - [Test Structure](#test-structure)
  - [Fake Data Generation](#fake-data-generation)
  - [HTTP Requests](#http-requests)
  - [Response Validation](#response-validation)
  - [Data Extraction](#data-extraction)
  - [Shell Commands](#shell-commands)
- [Configuration](#configuration)
- [Test Dependencies](#test-dependencies)
- [Real-World Examples](#real-world-examples)
- [CLI Reference](#cli-reference)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## 🎯 What is Business Driven Development?

Business Driven Development (BDD) for APIs means writing tests that:

1. **Describe behavior in business terms** - "Register a new user" instead of "POST /api/users"
2. **Focus on outcomes** - What should happen, not how it's implemented
3. **Tell a story** - Tests flow logically from one business scenario to the next
4. **Are readable by stakeholders** - Product managers can review and understand tests

### Traditional API Testing vs ABDD

**Traditional Testing:**
```javascript
test('POST /api/users returns 201', async () => {
  const response = await fetch('/api/users', {
    method: 'POST',
    body: JSON.stringify({ username: 'test', email: 'test@example.com' })
  });
  expect(response.status).toBe(201);
});
```

**ABDD Approach:**
```yaml
tests:
  - name: Register new user
    description: A visitor can create a new account with valid information
    fake:
      username: "{username}"
      email: "{email}"
    request:
      method: POST
      url: /users
      body: '{"username": "${username}", "email": "${email}"}'
    expect:
      status: 201
      json:
        message: "User created successfully"
    extract:
      - path: id
        as: user_id
```

## 🚀 Quick Start

Get up and running with ABDD in less than 5 minutes!

### 1. Install ABDD

```bash
go install github.com/davesavic/abdd@latest
```

### 2. Initialize Your Project

```bash
mkdir my-api-tests && cd my-api-tests
abdd init
```

This creates:
- `abdd.yaml` - Global configuration
- `tests/` directory with sample tests

### 3. Run Your First Test

```bash
abdd run --config abdd.yaml --folders tests
```

You'll see output like:
```
┌─────────────────────────────────┐
               Tests               
[1/3] ✓ Create post
[2/3] ✓ Create comment  
[3/3] ✓ Get comments for post

└─────────────────────────────────┘
┌─────────────────────────────────┐
              Summary              
Total: 3
Passed: 3
Pass rate: 100.0%
└─────────────────────────────────┘
```

That's it! You've just run your first business-driven API tests.

## 📦 Installation

### Option 1: Go Install (Recommended)

```bash
go install github.com/davesavic/abdd@latest
```

### Option 2: Download Binary

Download the latest release from [GitHub Releases](https://github.com/davesavic/abdd/releases):

```bash
# Linux/macOS
curl -L https://github.com/davesavic/abdd/releases/latest/download/abdd-linux-amd64 -o abdd
chmod +x abdd
sudo mv abdd /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/davesavic/abdd/releases/latest/download/abdd-windows-amd64.exe" -OutFile "abdd.exe"
```

### Option 3: Build from Source

```bash
git clone https://github.com/davesavic/abdd.git
cd abdd
go build -o abdd .
```

### Verify Installation

```bash
abdd --help
```

## 🧠 Core Concepts

### Test Files
- Tests are defined in YAML files (`.yaml` or `.yml`)
- Each file contains an array of tests under the `tests` key
- Tests run in dependency order, not file order

### Global Configuration
- `abdd.yaml` contains settings that apply to all tests
- Base URL, default headers, timeouts, etc.
- Can be overridden per test

### Variable Store
- ABDD maintains a key-value store during test execution
- Values can be extracted from responses or generated as fake data
- Variables are interpolated using `${variable_name}` syntax

### Test Lifecycle
Each test goes through these phases:
1. **Generate** fake data (if specified)
2. **Replace** variables in the test definition
3. **Execute** shell commands (if specified)  
4. **Make** HTTP request
5. **Validate** response
6. **Extract** data for future tests

## 📝 Writing Tests

### Test Structure

Every test follows this structure:

```yaml
tests:
  - name: "Human-readable test name"           # Required
    description: "Detailed explanation"        # Optional
    depends: ["Other Test Name"]               # Optional
    fake:                                      # Optional
      variable_name: "{faker_function}"
    command:                                   # Optional
      command: "echo 'Setting up data'"
      directory: "/path/to/work/dir"          # Optional
      as: "command_output"                     # Optional
    request:                                   # Required (unless command-only)
      method: POST
      url: /api/endpoint
      headers:                                 # Optional
        Authorization: "Bearer ${token}"
      body: '{"key": "${variable_name}"}'      # Optional
    expect:                                    # Required
      status: 200                             # Optional
      headers:                                # Optional
        Content-Type: "application/json"
      json:                                   # Optional
        field.path: "expected_value"
        nested.array.0.id: 123
    extract:                                  # Optional
      - path: "response.field"
        as: "variable_name"
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | **Required.** Unique identifier for the test |
| `description` | string | Human-readable explanation of what the test does |
| `depends` | array | Names of tests that must run before this one |
| `fake` | object | Map of variable names to faker functions |
| `command` | object | Shell command to execute before the request |
| `request` | object | HTTP request specification |
| `expect` | object | Response validation rules |
| `extract` | array | Values to extract from the response |

## 🎲 Fake Data Generation

ABDD integrates with [GoFakeIt](https://github.com/brianvoe/gofakeit) to provide **310+ functions** for generating realistic test data.

### Basic Syntax

```yaml
fake:
  variable_name: "{function_name}"
  with_params: "{function_name:param1,param2}"
```

### Complete Function Reference

#### 👤 Personal Information
```yaml
fake:
  name: "{name}"                    # "John Doe"
  first_name: "{firstname}"         # "John"
  last_name: "{lastname}"          # "Doe"
  middle_name: "{middlename}"      # "William"
  name_prefix: "{nameprefix}"       # "Mr."
  name_suffix: "{namesuffix}"       # "Jr."
  gender: "{gender}"                # "Male"
  ssn: "{ssn}"                      # "123-45-6789"
```

#### 📧 Contact Information
```yaml
fake:
  email: "{email}"                  # "john.doe@example.com"
  safe_email: "{safeemail}"         # "john.doe@example.org"
  username: "{username}"            # "johndoe"
  phone: "{phone}"                  # "(555) 123-4567"
  phone_formatted: "{phoneformatted}" # "+1-555-123-4567"
```

#### 🏢 Business & Professional
```yaml
fake:
  company: "{company}"              # "Acme Corp"
  company_suffix: "{companysuffix}" # "Inc"
  job_title: "{jobtitle}"          # "Software Engineer"
  job_descriptor: "{jobdescriptor}" # "Lead"
  job_level: "{joblevel}"          # "Manager"
  business_name: "{businessname}"   # "Smith & Associates"
  bs: "{bs}"                        # "synergize innovative solutions"
  buzzword: "{buzzword}"            # "paradigm"
```

#### 💰 Financial
```yaml
fake:
  currency_short: "{currencyshort}" # "USD"
  currency_long: "{currencylong}"   # "United States Dollar"
  price: "{price:10,100}"           # "45.99"
  credit_card: "{creditcardnumber}" # "4532623432341234"
  credit_card_cvv: "{creditcardcvv}" # "123"
  credit_card_exp: "{creditcardexp}" # "12/27"
  credit_card_type: "{creditcardtype}" # "Visa"
  achrouting: "{achrouting}"        # "123456789"
  achaccount: "{achaccount}"        # "123456789012"
  bitcoin_address: "{bitcoinaddress}" # "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
  bitcoin_private_key: "{bitcoinprivatekey}"
```

#### 🌍 Location & Geography  
```yaml
fake:
  address: "{address}"              # "123 Main St"
  street: "{street}"                # "Main St"
  street_number: "{streetnumber}"   # "123"
  city: "{city}"                    # "New York"
  state: "{state}"                  # "California"
  state_abr: "{stateabr}"          # "CA"
  zip: "{zip}"                      # "90210"
  country: "{country}"              # "United States"
  country_abr: "{countryabr}"      # "US"
  latitude: "{latitude}"            # "40.7128"
  longitude: "{longitude}"          # "-74.0060"
  timezone: "{timezone}"            # "America/New_York"
```

#### 🌐 Internet & Technology
```yaml
fake:
  domain_name: "{domainname}"       # "example.com"
  domain_suffix: "{domainsuffix}"   # "com"
  url: "{url}"                      # "https://example.com/path"
  ipv4_address: "{ipv4address}"     # "192.168.1.1"
  ipv6_address: "{ipv6address}"     # "2001:db8::1"
  mac_address: "{macaddress}"       # "01:23:45:67:89:ab"
  user_agent: "{useragent}"         # "Mozilla/5.0..."
  chrome_user_agent: "{chromeuseragent}"
  firefox_user_agent: "{firefoxuseragent}"
  safari_user_agent: "{safariuseragent}"
  http_method: "{httpmethod}"       # "GET"
  http_status_code: "{httpstatuscode}" # "200"
```

#### 📅 Date & Time
```yaml
fake:
  date: "{date}"                    # "2024-01-15"
  date_range: "{daterange:2024-01-01,2024-12-31}"
  past_date: "{pastdate}"           # Date in the past
  future_date: "{futuredate}"       # Date in the future
  past_time: "{pasttime}"           # Time in the past
  future_time: "{futuretime}"       # Time in the future
  month: "{month}"                  # "January"
  month_string: "{monthstring}"     # "Jan"
  day: "{day}"                      # "15"
  weekday: "{weekday}"              # "Monday"
  year: "{year}"                    # "2024"
  timezone_abv: "{timezoneabv}"     # "EST"
  timezone_full: "{timezonefull}"   # "Eastern Standard Time"
  nano_second: "{nanosecond}"       # "123456789"
```

#### 🎨 Visual & Creative
```yaml
fake:
  color: "{color}"                  # "Red"
  hex_color: "{hexcolor}"           # "#FF0000"
  rgb_color: "{rgbcolor}"           # "rgb(255, 0, 0)"
  safe_color: "{safecolor}"         # "Blue"
  image_url: "{imageurl:400,300}"   # "https://picsum.photos/400/300"
  image_jpeg: "{imagejpeg:400,300}"
  image_png: "{imagepng:400,300}"
```

#### 📱 Product & Commerce
```yaml
fake:
  product: "{product}"              # Complete product info
  product_name: "{productname}"     # "Wireless Headphones"
  product_description: "{productdescription}"
  product_category: "{productcategory}" # "Electronics"
  product_feature: "{productfeature}" # "Noise Cancelling"
  product_material: "{productmaterial}" # "Plastic"
  upc: "{upc}"                      # "123456789012"
  isbn: "{isbn}"                    # "978-3-16-148410-0"
```

#### 🎯 Text & Content
```yaml
fake:
  word: "{word}"                    # "fantastic"
  sentence: "{sentence:5}"          # 5-word sentence
  paragraph: "{paragraph:3,5,10,\n}" # 3-5 sentences, ~10 words each
  question: "{question}"            # "How are you today?"
  quote: "{quote}"                  # "Be yourself; everyone else is taken."
  phrase: "{phrase}"                # "all good things"
  lorem_ipsum_word: "{loremipsumword}" # "lorem"
  lorem_ipsum_sentence: "{loremipsumsentence:5}"
  lorem_ipsum_paragraph: "{loremipsumparagraph:3,5,10,\n}"
```

#### 🔢 Numbers & Identifiers
```yaml
fake:
  number: "{number:1,100}"          # Random number between 1-100
  int: "{int:1,100}"                # Integer between 1-100
  int8: "{int8}"                    # 8-bit integer
  int16: "{int16}"                  # 16-bit integer
  int32: "{int32}"                  # 32-bit integer  
  int64: "{int64}"                  # 64-bit integer
  uint: "{uint:1,100}"              # Unsigned int 1-100
  uint8: "{uint8}"                  # 8-bit unsigned int
  uint16: "{uint16}"                # 16-bit unsigned int
  uint32: "{uint32}"                # 32-bit unsigned int
  uint64: "{uint64}"                # 64-bit unsigned int
  float32: "{float32:1,100}"        # 32-bit float 1-100
  float64: "{float64:1,100}"        # 64-bit float 1-100
  uuid: "{uuid}"                    # "123e4567-e89b-12d3-a456-426614174000"
```

#### 🎲 Random Selections
```yaml
fake:
  random_string: "{randomstring:alphabet,10}" # Random 10-char string
  random_from_list: "{randomstring:[red,green,blue]}" # Pick from list
  bool: "{bool}"                    # "true" or "false"
  flip_a_coin: "{flipacoin}"        # "Heads" or "Tails"
  dice: "{dice:6}"                  # Roll 6-sided die (1-6)
```

#### 🍔 Food & Dining
```yaml
fake:
  food: "{food}"                    # "Pizza"
  fruit: "{fruit}"                  # "Apple"
  vegetable: "{vegetable}"          # "Carrot"
  breakfast: "{breakfast}"          # "Pancakes"
  lunch: "{lunch}"                  # "Sandwich"
  dinner: "{dinner}"                # "Steak"
  snack: "{snack}"                  # "Chips"
  dessert: "{dessert}"              # "Ice Cream"
  drink: "{drink}"                  # "Coffee"
```

#### 🍺 Beverages
```yaml
fake:
  beer_name: "{beername}"           # "Budweiser"
  beer_style: "{beerstyle}"         # "IPA"
  beer_hop: "{beerhop}"             # "Cascade"
  beer_yeast: "{beeryeast}"         # "Ale"
  beer_malt: "{beermalt}"           # "Pale"
  beer_ibu: "{beeribu}"             # "45"
  beer_alcohol: "{beeralcohol}"     # "5.2%"
  beer_blg: "{beerblg}"             # "12"
```

#### 🚗 Automotive
```yaml
fake:
  car: "{car}"                      # Complete car info
  car_maker: "{carmaker}"           # "Toyota"
  car_model: "{carmodel}"           # "Camry"
  car_type: "{cartype}"             # "Sedan"
  car_fuel_type: "{carfueltype}"    # "Gasoline"
  car_transmission_type: "{cartransmissiontype}" # "Automatic"
```

#### 🎮 Gaming
```yaml
fake:
  gamertag: "{gamertag}"            # "xXGamerXx"
  dice: "{dice:20}"                 # D&D 20-sided die
```

#### 🐱 Animals
```yaml
fake:
  animal: "{animal}"                # "Dog"
  animal_type: "{animaltype}"       # "Mammal"
  farm_animal: "{farmanimal}"       # "Cow"
  cat: "{cat}"                      # "Persian"
  dog: "{dog}"                      # "Golden Retriever"
  bird: "{bird}"                    # "Eagle"
```

#### 📚 Media & Entertainment
```yaml
fake:
  book: "{book}"                    # Complete book info
  book_title: "{booktitle}"         # "The Great Gatsby"
  book_author: "{bookauthor}"       # "F. Scott Fitzgerald"
  book_genre: "{bookgenre}"         # "Fiction"
  movie: "{movie}"                  # Complete movie info
  movie_name: "{moviename}"         # "The Matrix"
  movie_genre: "{moviegenre}"       # "Action"
```

#### 🏫 Education
```yaml
fake:
  school: "{school}"                # "Harvard University"
  college: "{college}"              # "MIT"
  university: "{university}"        # "Stanford University"
```

#### 🎵 Music
```yaml
fake:
  music: "{music}"                  # Complete song info
  music_artist: "{musicartist}"     # "The Beatles"
  music_album: "{musicalbum}"       # "Abbey Road"
  music_song: "{musicsong}"         # "Come Together"
  music_genre: "{musicgenre}"       # "Rock"
```

#### 💻 Technology
```yaml
fake:
  app: "{app}"                      # Complete app info
  app_name: "{appname}"             # "Instagram"
  app_version: "{appversion}"       # "2.1.0"
  programming_language: "{programminglanguage}" # "Go"
```

#### 🎭 Characters & Fiction
```yaml
fake:
  superhero: "{superhero}"          # "Superman"
  villain: "{villain}"              # "Joker"
  celebrity: "{celebrity}"          # "Leonardo DiCaprio"
```

### Advanced Fake Data Patterns

#### Parameters and Options
Many faker functions accept parameters:

```yaml
fake:
  # Numbers with range
  user_id: "{number:1000,9999}"     # Between 1000-9999
  
  # Text with length/count
  title: "{sentence:3}"             # 3-word sentence
  bio: "{paragraph:2,4,15,\n}"      # 2-4 sentences, ~15 words each
  
  # Dates with ranges  
  birth_date: "{daterange:1990-01-01,2000-12-31}"
  
  # Custom selections
  role: "{randomstring:[admin,user,moderator]}"
  priority: "{randomstring:[low,medium,high,critical]}"
```

#### Complex Combinations
```yaml
fake:
  # Realistic user data
  full_name: "{name}"
  email: "{safeemail}"
  username: "{username}"
  phone: "{phoneformatted}"
  birth_date: "{daterange:1980-01-01,2005-12-31}"
  
  # E-commerce product
  product_name: "{productname}"
  price: "{price:10,1000}"
  description: "{productdescription}"
  category: "{productcategory}"
  
  # Address info
  street_address: "{address}"
  city: "{city}"
  state: "{state}"
  postal_code: "{zip}"
  country: "{country}"
```

## 🌐 HTTP Requests

### Request Methods
ABDD supports all standard HTTP methods:

```yaml
request:
  method: GET     # GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
  url: /api/endpoint
```

### URLs and Paths
URLs are automatically prefixed with the `base_url` from your global config:

```yaml
# In abdd.yaml
global:
  config:
    base_url: https://api.example.com/v1

# In test file - becomes https://api.example.com/v1/users
request:
  url: /users
  
# Variable interpolation works too
request:
  url: /users/${user_id}/posts
```

### Headers
Set custom headers per request:

```yaml
request:
  headers:
    Authorization: "Bearer ${access_token}"
    Content-Type: "application/json"
    X-API-Version: "v1"
    User-Agent: "ABDD/1.0"
```

### Request Bodies
Send JSON, form data, or raw text:

#### JSON Bodies
```yaml
request:
  method: POST
  url: /users
  body: |
    {
      "username": "${username}",
      "email": "${email}",
      "profile": {
        "firstName": "${first_name}",
        "lastName": "${last_name}"
      }
    }
```

#### Form Data
```yaml
request:
  method: POST
  url: /login
  headers:
    Content-Type: "application/x-www-form-urlencoded"
  body: "username=${username}&password=${password}"
```

#### Raw Text
```yaml
request:
  method: POST
  url: /webhook
  headers:
    Content-Type: "text/plain"
  body: "Event: user.created\nUser ID: ${user_id}"
```

### Authentication Patterns

#### Bearer Token
```yaml
request:
  headers:
    Authorization: "Bearer ${jwt_token}"
```

#### Basic Auth
```yaml
request:
  headers:
    Authorization: "Basic ${base64_credentials}"
```

#### API Key
```yaml
request:
  headers:
    X-API-Key: "${api_key}"
```

#### Custom Auth
```yaml
request:
  headers:
    X-Auth-Token: "${auth_token}"
    X-User-ID: "${user_id}"
```

## ✅ Response Validation

ABDD provides flexible response validation using three main areas: status codes, headers, and JSON content.

### Status Code Validation

```yaml
expect:
  status: 200  # Exact status code match
```

Common status codes:
- `200` - OK
- `201` - Created  
- `204` - No Content
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `422` - Unprocessable Entity
- `500` - Internal Server Error

### Header Validation

Validate response headers:

```yaml
expect:
  headers:
    Content-Type: "application/json"
    Cache-Control: "no-cache"
    X-RateLimit-Remaining: "99"
```

### JSON Response Validation

ABDD uses [gjson](https://github.com/tidwall/gjson) for powerful JSON path validation:

#### Simple Fields
```yaml
expect:
  json:
    message: "User created successfully"
    success: true
    user.id: 123
    user.email: "john@example.com"
```

#### Nested Objects
```yaml
expect:
  json:
    user.profile.firstName: "John"
    user.profile.address.city: "New York"
    settings.preferences.theme: "dark"
```

#### Array Validation
```yaml
expect:
  json:
    users.#: 3                    # Array has 3 elements
    users.0.name: "John"          # First user's name
    users.-1.name: "Jane"         # Last user's name
    "tags.#.name": "important"    # Any tag has name "important"
```

#### Complex Queries
```yaml
expect:
  json:
    "users.#(age>18)#": 5         # 5 users over 18
    "products.#(price<100).name": "Budget Item"  # Product under $100
    "events.#(type==click)#": 10  # 10 click events
```

#### Type Validation
```yaml
expect:
  json:
    user.id: 123              # Number
    user.active: true         # Boolean
    user.name: "John"         # String
    user.metadata: null       # Null value
```

#### Advanced Pattern Matching
```yaml
expect:
  json:
    # Using wildcards
    "*.user_id": 123
    
    # Conditional matching
    "orders.#(status==shipped).total": 99.99
    
    # Multiple conditions
    "products.#(price>10 && category==electronics)#": 5
```

### Validation Examples by Use Case

#### API Authentication
```yaml
expect:
  status: 200
  headers:
    Content-Type: "application/json"
  json:
    token_type: "Bearer"
    expires_in: 3600
```

#### User Registration
```yaml
expect:
  status: 201
  json:
    message: "User created successfully" 
    user.id: "{{number}}"      # Any number
    user.created_at: "{{date}}" # Any date
```

#### Error Handling
```yaml
expect:
  status: 400
  json:
    error: "validation_failed"
    messages.#: 2              # Exactly 2 error messages
    messages.0: "Email is required"
```

#### Pagination
```yaml
expect:
  status: 200
  json:
    page: 1
    per_page: 10
    total: 100
    data.#: 10                 # 10 items on this page
```

## 📤 Data Extraction

Extract values from API responses to use in subsequent tests:

### Basic Extraction

```yaml
extract:
  - path: "id"                 # Extract response.id
    as: "user_id"             # Store as ${user_id}
  
  - path: "token"
    as: "auth_token"
    
  - path: "user.email"
    as: "user_email"
```

### Nested Data Extraction

```yaml
extract:
  - path: "user.profile.id"
    as: "profile_id"
    
  - path: "settings.api.key"
    as: "api_key"
    
  - path: "metadata.correlation_id"  
    as: "correlation_id"
```

### Array Element Extraction

```yaml
extract:
  - path: "items.0.id"         # First item's ID
    as: "first_item_id"
    
  - path: "users.-1.name"      # Last user's name  
    as: "last_user_name"
    
  - path: "tags.#.name"        # All tag names (as JSON array)
    as: "all_tag_names"
```

### Complex Query Extraction

```yaml
extract:
  # Find specific items
  - path: "products.#(featured==true).id"
    as: "featured_product_id"
    
  # Extract from matching condition
  - path: "orders.#(status==pending).total"
    as: "pending_order_total"
    
  # Count matching items
  - path: "events.#(type==error)#"
    as: "error_count"
```

### Data Type Handling

ABDD automatically handles different data types:

```yaml
# These are extracted with proper types
extract:
  - path: "user.id"           # Number → stored as number
    as: "user_id"
    
  - path: "user.active"       # Boolean → stored as boolean  
    as: "is_active"
    
  - path: "user.name"         # String → stored as string
    as: "username"
    
  - path: "user.tags"         # Array → stored as JSON array
    as: "user_tags"
    
  - path: "user.metadata"     # Object → stored as JSON object  
    as: "user_metadata"
```

### Using Extracted Variables

Once extracted, variables can be used in subsequent tests:

```yaml
tests:
  - name: Create user
    request:
      method: POST
      url: /users
      body: '{"name": "John"}'
    extract:
      - path: "id" 
        as: "user_id"
        
  - name: Get user profile
    depends: ["Create user"]
    request:
      method: GET
      url: /users/${user_id}/profile    # Uses extracted user_id
      
  - name: Update user settings
    depends: ["Create user"] 
    request:
      method: PUT
      url: /users/${user_id}/settings
      body: '{"theme": "dark"}'
```

### Extraction Error Handling

If an extraction path doesn't exist, the test will fail:

```yaml
# This will fail if "nonexistent" field is missing
extract:
  - path: "nonexistent.field"
    as: "will_fail"
```

To make extraction optional, ensure the field exists first:

```yaml
# Better: validate the field exists first
expect:
  json:
    optional_field: "{{any}}"  # Ensures field exists
extract:
  - path: "optional_field"
    as: "extracted_value"
```

## 🖥️ Shell Commands

Execute shell commands as part of your test flow:

### Basic Command Execution

```yaml
command:
  command: "echo 'Setting up test data'"
```

### Commands with Working Directory

```yaml
command:
  command: "npm install"
  directory: "/path/to/project"
```

### Capturing Command Output

```yaml
command:
  command: "cat config.json | jq .database.host"
  as: "db_host"              # Store output as ${db_host}
```

### Variable Interpolation in Commands

```yaml
command:
  command: "curl -X POST https://api.example.com/users -d '{\"name\": \"${user_name}\"}'"
```

### Common Use Cases

#### Database Seeding
```yaml
command:
  command: "mysql -u root -p${db_password} myapp < seed.sql"
  directory: "/path/to/sql/files"
```

#### File Manipulation  
```yaml
command:
  command: "jq '.test_mode = true' config.json > temp_config.json"
  as: "config_updated"
```

#### Environment Setup
```yaml
command:
  command: "docker-compose up -d database"
  directory: "/path/to/docker"
```

#### API Preparation
```yaml
command:
  command: "kubectl apply -f test-resources.yaml"
  as: "k8s_output"
```

### Command-Only Tests

Some tests might only run commands without HTTP requests:

```yaml
tests:
  - name: Setup test database
    description: Create and seed test database
    command:
      command: "make setup-test-db"
      directory: "/app"
    expect: {}  # Still need expect, but can be empty
```

### Error Handling

If a command fails (non-zero exit code), the test fails:

```yaml
command:
  command: "exit 1"  # This will fail the test
```

## ⚙️ Configuration

### Global Configuration (abdd.yaml)

```yaml
global:
  config:
    base_url: "https://api.example.com/v1"
    headers:
      Content-Type: "application/json"
      User-Agent: "ABDD/1.0"
    timeout: 30                # Request timeout in seconds
    stop_on_error: true        # Stop execution on first failure
    verbose: false             # Enable detailed output
```

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `base_url` | string | - | Base URL prepended to all request URLs |
| `headers` | object | `{}` | Default headers added to all requests |
| `timeout` | integer | `30` | HTTP request timeout in seconds |
| `stop_on_error` | boolean | `true` | Stop test execution on first failure |
| `verbose` | boolean | `false` | Enable verbose output during execution |

### Environment-Specific Configurations

Create different config files for different environments:

```bash
# Development
abdd run --config abdd.dev.yaml --folders tests

# Staging  
abdd run --config abdd.staging.yaml --folders tests

# Production
abdd run --config abdd.prod.yaml --folders tests
```

Example environment configs:

```yaml
# abdd.dev.yaml
global:
  config:
    base_url: "http://localhost:3000/api/v1"
    timeout: 60
    verbose: true

# abdd.staging.yaml  
global:
  config:
    base_url: "https://staging-api.example.com/v1"
    timeout: 30
    headers:
      X-Environment: "staging"

# abdd.prod.yaml
global:
  config:
    base_url: "https://api.example.com/v1"
    timeout: 10
    stop_on_error: false
    headers:
      X-Environment: "production"
```

## 🔗 Test Dependencies

### Declaring Dependencies

Tests can depend on other tests and will run in the correct order:

```yaml
tests:
  - name: Setup user account
    # ... test definition
    
  - name: Login user
    depends: ["Setup user account"]
    # ... test definition
    
  - name: Create user profile  
    depends: ["Login user"]
    # ... test definition
    
  - name: Update profile settings
    depends: ["Create user profile", "Login user"]  # Multiple dependencies
    # ... test definition
```

### Dependency Resolution

ABDD automatically:
1. **Sorts tests** based on dependencies using topological sorting
2. **Detects circular dependencies** and reports errors
3. **Validates dependencies** exist before execution
4. **Shares variable store** across dependent tests

### Complex Dependency Chains

```yaml
tests:
  # Base setup
  - name: Create admin user
    # ... creates admin account
    
  - name: Create regular user  
    # ... creates regular account
    
  # Admin operations (depend on admin user)
  - name: Create organization
    depends: ["Create admin user"]
    # ... admin creates org
    
  - name: Set organization settings
    depends: ["Create organization"]
    # ... admin configures org
    
  # User operations (depend on both users and org)
  - name: Join organization
    depends: ["Create regular user", "Create organization"]
    # ... regular user joins org
    
  - name: Access organization data
    depends: ["Join organization", "Set organization settings"]
    # ... user accesses configured org data
```

### Cross-File Dependencies

Dependencies work across multiple test files:

```yaml
# tests/auth.yaml
tests:
  - name: User registration
    # ... registration test
    
# tests/profile.yaml  
tests:
  - name: Create user profile
    depends: ["User registration"]  # References test from auth.yaml
    # ... profile creation
```

### Dependency Best Practices

#### 1. Keep Dependencies Minimal
```yaml
# Good: Only depend on what you need
- name: Update user email
  depends: ["Create user"]
  
# Avoid: Over-depending  
- name: Update user email
  depends: ["Create user", "Login user", "Load user settings"]
```

#### 2. Use Descriptive Test Names
```yaml
# Good: Clear, specific names
depends: ["Create admin user account", "Setup organization permissions"]

# Avoid: Generic names
depends: ["Test1", "Setup", "User stuff"]
```

#### 3. Group Related Tests
```yaml
# tests/user-lifecycle.yaml
tests:
  - name: Register new user
  - name: Verify user email
    depends: ["Register new user"]
  - name: Complete user profile  
    depends: ["Verify user email"]
  - name: Activate user account
    depends: ["Complete user profile"]
```

### Error Handling with Dependencies

If a dependency fails:
- **Dependent tests are skipped** (not failed)
- **Error message shows** which dependency failed
- **Execution continues** with independent tests (unless `stop_on_error: true`)

Example output:
```
[1/4] ✗ Create user account
       → Error: User registration failed
[2/4] ⊘ Update user profile (skipped - dependency failed)
[3/4] ⊘ Delete user account (skipped - dependency failed)  
[4/4] ✓ Get system status
```

## 🌟 Real-World Examples

### E-Commerce API Testing

Complete workflow testing user registration through order completion:

```yaml
# tests/ecommerce-flow.yaml
tests:
  - name: Register new customer
    description: Customer creates account with email and password
    fake:
      customer_email: "{email}"
      customer_password: "{password:8,16}"
      first_name: "{firstname}"
      last_name: "{lastname}"
    request:
      method: POST
      url: /customers
      body: |
        {
          "email": "${customer_email}",
          "password": "${customer_password}",
          "firstName": "${first_name}",
          "lastName": "${last_name}"
        }
    expect:
      status: 201
      json:
        message: "Customer registered successfully"
        customer.id: "{{number}}"
        customer.email: "${customer_email}"
    extract:
      - path: "customer.id"
        as: "customer_id"
      - path: "access_token"
        as: "auth_token"

  - name: Browse product catalog
    description: Get list of available products
    depends: ["Register new customer"]
    request:
      method: GET
      url: /products
      headers:
        Authorization: "Bearer ${auth_token}"
    expect:
      status: 200
      json:
        products.#: "{{gt:0}}"  # Has at least one product
    extract:
      - path: "products.0.id"
        as: "product_id"
      - path: "products.0.price"
        as: "product_price"

  - name: Add item to cart
    description: Customer adds product to shopping cart
    depends: ["Browse product catalog"]
    fake:
      quantity: "{number:1,3}"
    request:
      method: POST
      url: /customers/${customer_id}/cart/items
      headers:
        Authorization: "Bearer ${auth_token}"
      body: |
        {
          "productId": "${product_id}",
          "quantity": ${quantity}
        }
    expect:
      status: 200
      json:
        cart.items.#: 1
        cart.items.0.productId: "${product_id}"
        cart.total: "{{number}}"
    extract:
      - path: "cart.id"
        as: "cart_id"
      - path: "cart.total"
        as: "cart_total"

  - name: Create shipping address
    description: Customer adds shipping address
    depends: ["Register new customer"]
    fake:
      street: "{street}"
      city: "{city}"
      state: "{stateabr}"
      zip: "{zip}"
      country: "{countryabr}"
    request:
      method: POST
      url: /customers/${customer_id}/addresses
      headers:
        Authorization: "Bearer ${auth_token}"
      body: |
        {
          "type": "shipping",
          "street": "${street}",
          "city": "${city}",
          "state": "${state}",
          "zipCode": "${zip}",
          "country": "${country}"
        }
    expect:
      status: 201
      json:
        address.id: "{{number}}"
        address.type: "shipping"
    extract:
      - path: "address.id"
        as: "shipping_address_id"

  - name: Process payment
    description: Customer pays for items in cart
    depends: ["Add item to cart", "Create shipping address"]
    fake:
      card_number: "{creditcardnumber}"
      card_cvv: "{creditcardcvv}"
      card_exp: "{creditcardexp}"
    request:
      method: POST
      url: /payments
      headers:
        Authorization: "Bearer ${auth_token}"
      body: |
        {
          "cartId": "${cart_id}",
          "paymentMethod": {
            "type": "credit_card",
            "cardNumber": "${card_number}",
            "cvv": "${card_cvv}",
            "expiryDate": "${card_exp}"
          },
          "billingAddress": {
            "addressId": "${shipping_address_id}"
          }
        }
    expect:
      status: 200
      json:
        payment.status: "completed"
        payment.amount: "${cart_total}"
    extract:
      - path: "payment.id"
        as: "payment_id"

  - name: Create order
    description: Convert paid cart to order
    depends: ["Process payment"]
    request:
      method: POST
      url: /orders
      headers:
        Authorization: "Bearer ${auth_token}"
      body: |
        {
          "cartId": "${cart_id}",
          "paymentId": "${payment_id}",
          "shippingAddressId": "${shipping_address_id}"
        }
    expect:
      status: 201
      json:
        order.id: "{{number}}"
        order.status: "confirmed"
        order.customerId: "${customer_id}"
        order.total: "${cart_total}"
    extract:
      - path: "order.id"
        as: "order_id"
      - path: "order.trackingNumber"
        as: "tracking_number"

  - name: Check order status
    description: Customer checks order status
    depends: ["Create order"]
    request:
      method: GET
      url: /orders/${order_id}
      headers:
        Authorization: "Bearer ${auth_token}"
    expect:
      status: 200
      json:
        order.id: "${order_id}"
        order.status: "confirmed"
        order.trackingNumber: "${tracking_number}"
```

### Multi-Tenant SaaS API Testing

Testing tenant management and isolation:

```yaml
# tests/saas-multitenancy.yaml
tests:
  - name: Create tenant organization
    description: Super admin creates new tenant organization
    fake:
      org_name: "{company}"
      admin_email: "{email}"
      admin_password: "{password:10,20}"
    command:
      command: "echo 'Testing tenant: ${org_name}'"
    request:
      method: POST
      url: /admin/tenants
      headers:
        Authorization: "Bearer ${super_admin_token}"
      body: |
        {
          "name": "${org_name}",
          "adminUser": {
            "email": "${admin_email}",
            "password": "${admin_password}"
          }
        }
    expect:
      status: 201
      json:
        tenant.id: "{{uuid}}"
        tenant.name: "${org_name}"
        tenant.status: "active"
    extract:
      - path: "tenant.id"
        as: "tenant_id"
      - path: "tenant.apiKey"
        as: "tenant_api_key"
      - path: "adminUser.id"
        as: "tenant_admin_id"

  - name: Login as tenant admin
    description: Tenant admin authenticates
    depends: ["Create tenant organization"]
    request:
      method: POST
      url: /auth/login
      body: |
        {
          "email": "${admin_email}",
          "password": "${admin_password}",
          "tenantId": "${tenant_id}"
        }
    expect:
      status: 200
      json:
        token_type: "Bearer"
        expires_in: 3600
        user.tenantId: "${tenant_id}"
    extract:
      - path: "access_token"
        as: "tenant_admin_token"

  - name: Create tenant user
    description: Tenant admin creates regular user
    depends: ["Login as tenant admin"]
    fake:
      user_email: "{email}"
      user_name: "{name}"
    request:
      method: POST
      url: /users
      headers:
        Authorization: "Bearer ${tenant_admin_token}"
        X-Tenant-ID: "${tenant_id}"
      body: |
        {
          "email": "${user_email}",
          "name": "${user_name}",
          "role": "user"
        }
    expect:
      status: 201
      json:
        user.tenantId: "${tenant_id}"
        user.email: "${user_email}"
        user.role: "user"
    extract:
      - path: "user.id"
        as: "tenant_user_id"

  - name: Test tenant isolation
    description: Verify tenant cannot access other tenant's data
    depends: ["Create tenant user"]
    request:
      method: GET
      url: /users
      headers:
        Authorization: "Bearer ${tenant_admin_token}"
        X-Tenant-ID: "different-tenant-id"  # Try to access different tenant
    expect:
      status: 403
      json:
        error: "access_denied"
        message: "Insufficient permissions for tenant"

  - name: Create tenant-specific resource
    description: Create resource within tenant scope
    depends: ["Create tenant user"]
    fake:
      project_name: "{bs}"
      project_description: "{sentence:10}"
    request:
      method: POST
      url: /projects
      headers:
        Authorization: "Bearer ${tenant_admin_token}"
        X-Tenant-ID: "${tenant_id}"
      body: |
        {
          "name": "${project_name}",
          "description": "${project_description}",
          "ownerId": "${tenant_user_id}"
        }
    expect:
      status: 201
      json:
        project.tenantId: "${tenant_id}"
        project.ownerId: "${tenant_user_id}"
    extract:
      - path: "project.id"
        as: "project_id"

  - name: Verify tenant data isolation
    description: List tenant projects and verify isolation
    depends: ["Create tenant-specific resource"]
    request:
      method: GET
      url: /projects
      headers:
        Authorization: "Bearer ${tenant_admin_token}"
        X-Tenant-ID: "${tenant_id}"
    expect:
      status: 200
      json:
        projects.#: 1                        # Only one project
        projects.0.id: "${project_id}"
        projects.0.tenantId: "${tenant_id}"  # Belongs to correct tenant
```

### Microservices Integration Testing

Testing communication between multiple services:

```yaml
# tests/microservices-integration.yaml
tests:
  - name: Health check all services
    description: Verify all microservices are running
    request:
      method: GET
      url: /health
    expect:
      status: 200
      json:
        services.user-service: "healthy"
        services.order-service: "healthy"
        services.inventory-service: "healthy"
        services.payment-service: "healthy"

  - name: Create user in user service
    description: User service creates new user
    depends: ["Health check all services"]
    fake:
      email: "{email}"
      username: "{username}"
    request:
      method: POST
      url: /user-service/users
      body: |
        {
          "email": "${email}",
          "username": "${username}"
        }
    expect:
      status: 201
      json:
        user.id: "{{uuid}}"
    extract:
      - path: "user.id"
        as: "user_id"

  - name: Check inventory levels
    description: Inventory service shows available products
    depends: ["Health check all services"]
    request:
      method: GET
      url: /inventory-service/products
    expect:
      status: 200
      json:
        products.#: "{{gt:0}}"
        "products.#(stock>0)#": "{{gt:0}}"  # At least one in stock
    extract:
      - path: "products.#(stock>0).id"
        as: "available_product_id"
      - path: "products.#(stock>0).price"
        as: "product_price"

  - name: Create order in order service
    description: Order service creates order and calls inventory
    depends: ["Create user in user service", "Check inventory levels"]
    fake:
      quantity: "{number:1,2}"
    request:
      method: POST
      url: /order-service/orders
      body: |
        {
          "userId": "${user_id}",
          "items": [
            {
              "productId": "${available_product_id}",
              "quantity": ${quantity}
            }
          ]
        }
    expect:
      status: 201
      json:
        order.userId: "${user_id}"
        order.status: "pending"
        order.total: "{{number}}"
    extract:
      - path: "order.id"
        as: "order_id"
      - path: "order.total"
        as: "order_total"

  - name: Process payment via payment service
    description: Payment service processes order payment
    depends: ["Create order in order service"]
    fake:
      payment_method: "{randomstring:[credit_card,debit_card,paypal]}"
    request:
      method: POST
      url: /payment-service/payments
      body: |
        {
          "orderId": "${order_id}",
          "amount": ${order_total},
          "paymentMethod": "${payment_method}",
          "userId": "${user_id}"
        }
    expect:
      status: 200
      json:
        payment.status: "completed"
        payment.orderId: "${order_id}"
    extract:
      - path: "payment.id"
        as: "payment_id"

  - name: Verify order completion
    description: Order service updates status after payment
    depends: ["Process payment via payment service"]
    command:
      command: "sleep 2"  # Wait for async processing
    request:
      method: GET
      url: /order-service/orders/${order_id}
    expect:
      status: 200
      json:
        order.status: "completed"
        order.paymentId: "${payment_id}"

  - name: Verify inventory reduction
    description: Inventory service reduced stock after order
    depends: ["Verify order completion"]
    request:
      method: GET
      url: /inventory-service/products/${available_product_id}
    expect:
      status: 200
      json:
        product.id: "${available_product_id}"
        # Note: Can't easily test exact stock reduction without knowing initial value
        # In real scenarios, you might extract initial stock first
```

### Authentication & Authorization Flow

Complete auth flow testing:

```yaml
# tests/auth-flow.yaml
tests:
  - name: Register with weak password
    description: Registration should fail with weak password
    fake:
      email: "{email}"
      weak_password: "123"  # Intentionally weak
    request:
      method: POST
      url: /auth/register
      body: |
        {
          "email": "${email}",
          "password": "${weak_password}"
        }
    expect:
      status: 422
      json:
        error: "validation_error"
        "errors.password": "Password must be at least 8 characters"

  - name: Register with valid credentials
    description: User registration with strong password
    fake:
      email: "{email}"
      strong_password: "{password:12,20}"
      first_name: "{firstname}"
      last_name: "{lastname}"
    request:
      method: POST
      url: /auth/register
      body: |
        {
          "email": "${email}",
          "password": "${strong_password}",
          "firstName": "${first_name}",
          "lastName": "${last_name}"
        }
    expect:
      status: 201
      json:
        message: "Registration successful"
        user.email: "${email}"
    extract:
      - path: "user.id"
        as: "user_id"

  - name: Login with invalid credentials
    description: Login should fail with wrong password
    depends: ["Register with valid credentials"]
    request:
      method: POST
      url: /auth/login
      body: |
        {
          "email": "${email}",
          "password": "wrong_password"
        }
    expect:
      status: 401
      json:
        error: "invalid_credentials"

  - name: Login with valid credentials
    description: User logs in successfully
    depends: ["Register with valid credentials"]
    request:
      method: POST
      url: /auth/login
      body: |
        {
          "email": "${email}",
          "password": "${strong_password}"
        }
    expect:
      status: 200
      json:
        token_type: "Bearer"
        expires_in: 3600
        user.id: "${user_id}"
    extract:
      - path: "access_token"
        as: "access_token"
      - path: "refresh_token"
        as: "refresh_token"

  - name: Access protected resource
    description: Use token to access protected endpoint
    depends: ["Login with valid credentials"]
    request:
      method: GET
      url: /profile
      headers:
        Authorization: "Bearer ${access_token}"
    expect:
      status: 200
      json:
        user.id: "${user_id}"
        user.email: "${email}"

  - name: Access protected resource without token
    description: Should fail without authentication
    request:
      method: GET
      url: /profile
    expect:
      status: 401
      json:
        error: "unauthorized"

  - name: Refresh access token
    description: Get new access token using refresh token
    depends: ["Login with valid credentials"]
    command:
      command: "sleep 1"  # Simulate token near expiry
    request:
      method: POST
      url: /auth/refresh
      body: |
        {
          "refreshToken": "${refresh_token}"
        }
    expect:
      status: 200
      json:
        token_type: "Bearer"
        expires_in: 3600
    extract:
      - path: "access_token"
        as: "new_access_token"

  - name: Use refreshed token
    description: Verify new token works
    depends: ["Refresh access token"]
    request:
      method: GET
      url: /profile
      headers:
        Authorization: "Bearer ${new_access_token}"
    expect:
      status: 200
      json:
        user.id: "${user_id}"

  - name: Logout user
    description: Invalidate user session
    depends: ["Use refreshed token"]
    request:
      method: POST
      url: /auth/logout
      headers:
        Authorization: "Bearer ${new_access_token}"
      body: |
        {
          "refreshToken": "${refresh_token}"
        }
    expect:
      status: 200
      json:
        message: "Logout successful"

  - name: Verify token invalidated
    description: Old token should no longer work
    depends: ["Logout user"]
    request:
      method: GET
      url: /profile
      headers:
        Authorization: "Bearer ${new_access_token}"
    expect:
      status: 401
      json:
        error: "invalid_token"
```

## 📖 CLI Reference

### Commands

#### `abdd init`
Initialize a new ABDD project with sample configuration and tests.

```bash
abdd init
```

**Creates:**
- `abdd.yaml` - Global configuration file
- `tests/` directory
- `tests/posts.yaml` - Sample test file with examples

#### `abdd run`
Execute test suites from specified folders.

```bash
abdd run --config abdd.yaml --folders tests
abdd run --config abdd.yaml --folders tests,integration --verbose
```

**Options:**
- `--config, -c` - Configuration file path (required)
- `--folders, -f` - Comma-separated list of test folders (required)  
- `--verbose, -v` - Enable verbose output

**Examples:**
```bash
# Basic execution
abdd run --config abdd.yaml --folders tests

# Multiple test folders
abdd run --config abdd.yaml --folders unit,integration,e2e

# Verbose output
abdd run --config abdd.yaml --folders tests --verbose

# Different environment
abdd run --config abdd.staging.yaml --folders tests
```

#### `abdd --help`
Show help information and available commands.

### Global Flags

These flags work with any command:

- `--config` - Specify configuration file path
- `--help, -h` - Show help for any command

### Exit Codes

- `0` - All tests passed
- `1` - One or more tests failed or execution error

### Verbose Output

With `--verbose` flag, ABDD shows detailed execution steps:

```
┌─────────────────────────────────┐
               Tests               

[1/3] Create post
  → Generate fake data: title="Lorem ipsum", body="Dolor sit amet..."
  → Replace variables in request body
  → Execute command: echo 'Creating post with title: Lorem ipsum'
  → Make HTTP request: POST https://api.example.com/posts
  → Validate response: status=201 ✓, headers ✓, json ✓
  → Extract data: postId=123
  ✓ Create post

[2/3] Create comment
  → Generate fake data: email="user@example.com", name="John Doe"
  → Replace variables: postId=123
  → Make HTTP request: POST https://api.example.com/comments  
  → Validate response: status=201 ✓
  → Extract data: commentId=456
  ✓ Create comment

[3/3] Get comments for post
  → Make HTTP request: GET https://api.example.com/posts/1/comments
  → Validate response: status=200 ✓, json.#=5 ✓, json.0.postId=1 ✓
  ✓ Get comments for post

└─────────────────────────────────┘
```

## 🎯 Best Practices

### Test Organization

#### File Structure
```
project/
├── abdd.yaml                 # Global config
├── tests/
│   ├── auth/                 # Authentication tests
│   │   ├── registration.yaml
│   │   ├── login.yaml
│   │   └── password-reset.yaml
│   ├── user/                 # User management tests
│   │   ├── profile.yaml
│   │   ├── preferences.yaml
│   │   └── permissions.yaml
│   ├── api/                  # Core API tests
│   │   ├── products.yaml
│   │   ├── orders.yaml
│   │   └── payments.yaml
│   └── integration/          # End-to-end flows
│       ├── checkout-flow.yaml
│       └── admin-workflow.yaml
└── configs/                  # Environment configs
    ├── abdd.dev.yaml
    ├── abdd.staging.yaml
    └── abdd.prod.yaml
```

#### Naming Conventions
```yaml
# Good: Descriptive, business-focused names
tests:
  - name: "Customer can register with valid email"
  - name: "Admin can create product with inventory"
  - name: "Payment fails with expired credit card"

# Avoid: Technical, generic names  
tests:
  - name: "POST /users returns 201"
  - name: "Test 1"
  - name: "API call"
```

### Variable Management

#### Use Meaningful Names
```yaml
# Good: Clear variable names
fake:
  customer_email: "{email}"
  product_price: "{price:10,100}"
  order_quantity: "{number:1,5}"

# Avoid: Generic names
fake:
  var1: "{email}"
  data: "{price:10,100}"
  x: "{number:1,5}"
```

#### Group Related Variables
```yaml
# Good: Grouped user data
fake:
  user_first_name: "{firstname}"
  user_last_name: "{lastname}" 
  user_email: "{email}"
  user_phone: "{phone}"

# Good: Grouped product data
fake:
  product_name: "{productname}"
  product_category: "{productcategory}"
  product_price: "{price:10,500}"
  product_description: "{productdescription}"
```

### Error Handling

#### Test Both Success and Failure Cases
```yaml
tests:
  - name: "Create user with valid data succeeds"
    # ... success case
    expect:
      status: 201

  - name: "Create user with duplicate email fails"
    # ... failure case  
    expect:
      status: 409
      json:
        error: "email_already_exists"

  - name: "Create user without required fields fails"
    # ... validation failure
    expect:
      status: 422
      json:
        error: "validation_error"
```

#### Validate Error Messages
```yaml
expect:
  status: 400
  json:
    error: "validation_failed"
    message: "Email is required"
    field: "email"
```

### Configuration Management

#### Environment-Specific Settings
```yaml
# abdd.dev.yaml - Development
global:
  config:
    base_url: "http://localhost:3000/api"
    timeout: 60
    verbose: true
    stop_on_error: false

# abdd.prod.yaml - Production  
global:
  config:
    base_url: "https://api.example.com"
    timeout: 10
    verbose: false
    stop_on_error: true
    headers:
      X-Environment: "production"
```

#### Shared Headers
```yaml
global:
  config:
    headers:
      Content-Type: "application/json"
      Accept: "application/json"
      User-Agent: "ABDD Test Suite"
      X-API-Version: "v1"
```

### Performance Considerations

#### Minimize Dependencies
```yaml
# Good: Minimal dependencies
tests:
  - name: "Get user profile"
    depends: ["Login user"]

# Avoid: Unnecessary dependencies
tests:  
  - name: "Get user profile"
    depends: ["Register user", "Verify email", "Login user", "Update settings"]
```

#### Use Realistic Test Data
```yaml
# Good: Realistic data sizes
fake:
  description: "{sentence:10}"        # ~10 words
  bio: "{paragraph:2,3,20,\n}"       # 2-3 sentences, ~20 words each

# Avoid: Excessive data
fake:
  description: "{sentence:100}"       # Unrealistically long
  bio: "{paragraph:20,30,50,\n}"     # Too verbose
```

### Maintainability

#### Document Complex Logic
```yaml
tests:
  - name: "Calculate shipping cost for international order"
    description: |
      Tests shipping calculation for international orders where:
      - Base shipping = $15
      - Each additional item = +$3
      - International surcharge = +$10
      - Orders over $100 get free upgrade to express
    # ... test implementation
```

#### Use Consistent Patterns
```yaml
# Consistent auth pattern across tests
request:
  headers:
    Authorization: "Bearer ${access_token}"

# Consistent error validation pattern  
expect:
  status: 400
  json:
    error: "{{string}}"
    message: "{{string}}"
```

### CI/CD Integration

#### Separate Test Suites
```bash
# Fast smoke tests
abdd run --config abdd.yaml --folders smoke

# Full regression suite  
abdd run --config abdd.yaml --folders tests

# Critical path only
abdd run --config abdd.yaml --folders critical
```

#### Environment-Aware Execution
```bash
# Different configs per environment
abdd run --config abdd.${ENVIRONMENT}.yaml --folders tests

# Skip certain tests in production
if [ "$ENVIRONMENT" != "production" ]; then
  abdd run --config abdd.yaml --folders destructive-tests
fi
```

## 🔧 Troubleshooting

### Common Issues

#### Configuration Problems

**Problem:** `Error: config file abdd.yaml does not exist`
```bash
# Solution: Initialize project first
abdd init

# Or specify correct path
abdd run --config path/to/abdd.yaml --folders tests
```

**Problem:** `Error: no folders provided`
```bash
# Solution: Specify test folders
abdd run --config abdd.yaml --folders tests

# Multiple folders
abdd run --config abdd.yaml --folders unit,integration
```

#### Test Execution Issues

**Problem:** `Error: circular dependency detected involving test 'Test A'`
```yaml
# Problem: Tests depend on each other
tests:
  - name: "Test A"
    depends: ["Test B"]
  - name: "Test B"  
    depends: ["Test A"]

# Solution: Remove circular dependency
tests:
  - name: "Setup data"
  - name: "Test A"
    depends: ["Setup data"]
  - name: "Test B"
    depends: ["Setup data"]
```

**Problem:** `Error: test 'Login user' depends on non-existent test 'Register user'`
```yaml
# Solution: Check test name spelling
depends: ["Register new user"]  # Make sure this matches exactly

# Or add missing test
tests:
  - name: "Register new user"  # This test must exist
    # ... 
  - name: "Login user"
    depends: ["Register new user"]
```

#### HTTP Request Issues

**Problem:** `Error: failed to make request: connection refused`
```yaml
# Check base_url is correct
global:
  config:
    base_url: "http://localhost:3000"  # Is service running on this port?
```

**Problem:** `Error: unexpected status code: expected 200, got 404`
```yaml
# Check URL path
request:
  url: /api/users  # Is this the correct endpoint?
  
# Enable verbose mode to see full request
abdd run --config abdd.yaml --folders tests --verbose
```

#### Validation Errors

**Problem:** `Error: json path not found: expected user.id to be present`
```yaml
# Check JSON response structure
# Enable verbose mode to see actual response
abdd run --config abdd.yaml --folders tests --verbose

# Or adjust JSON path
expect:
  json:
    id: "{{number}}"      # Instead of user.id
```

**Problem:** `Error: header not found: expected Content-Type to be present`
```yaml
# Check actual response headers (case sensitive)
expect:
  headers:
    content-type: "application/json"  # Try lowercase
```

#### Variable Issues

**Problem:** Variables not being replaced: `${user_id}` appears literally in request
```yaml
# Check variable was extracted
extract:
  - path: "id"
    as: "user_id"    # Variable name must match exactly

# Check test dependencies  
depends: ["Create user"]  # Test that extracts user_id must run first
```

**Problem:** `Error: failed to generate fake data for email`
```yaml
# Check faker function name
fake:
  email: "{email}"      # Correct
  # email: "{emaill}"   # Typo would cause error
```

### Debugging Tips

#### Enable Verbose Output
```bash
abdd run --config abdd.yaml --folders tests --verbose
```

This shows:
- Generated fake data values
- Variable substitutions  
- Full HTTP requests and responses
- Validation steps
- Extraction results

#### Isolate Problem Tests
```yaml
# Run single test file
abdd run --config abdd.yaml --folders specific-test-folder

# Create minimal test to debug issue
tests:
  - name: "Debug test"
    request:
      method: GET
      url: /health
    expect:
      status: 200
```

#### Check Service Health
```yaml
# Add health check test
tests:
  - name: "Service health check"
    request:
      method: GET  
      url: /health
    expect:
      status: 200
```

#### Validate JSON Paths
Use a JSON path tester with your response data:
- [JSONPath Online Evaluator](https://jsonpath.com/)
- [gjson Playground](https://gjson.dev/)

```yaml
# Test JSON paths separately
expect:
  json:
    user: "{{any}}"        # Check if user object exists
    user.id: "{{any}}"     # Check if user.id exists  
    user.id: "{{number}}"  # Check if user.id is number
```

#### Network Issues
```bash
# Test connectivity manually
curl -v http://localhost:3000/api/health

# Check if service is running
netstat -tulpn | grep :3000
```

### Getting Help

#### Check Logs
Enable verbose output and check error details:
```bash
abdd run --config abdd.yaml --folders tests --verbose 2>&1 | tee debug.log
```

#### Minimal Reproduction
Create a minimal test that reproduces the issue:
```yaml
tests:
  - name: "Minimal repro"
    request:
      method: GET
      url: /simple-endpoint
    expect:
      status: 200
```

#### Common Error Messages

| Error | Meaning | Solution |
|-------|---------|----------|
| `connection refused` | Service not running | Start your API service |
| `unexpected status code` | Wrong HTTP status | Check endpoint behavior |
| `json path not found` | Invalid JSON path | Verify response structure |
| `circular dependency` | Tests depend on each other | Fix dependency chain |
| `extraction path not found` | JSON path doesn't exist | Check response format |
| `failed to generate fake data` | Invalid faker function | Check function name |

## 🤝 Contributing

We welcome contributions to make ABDD even better! Here's how you can help:

### Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/davesavic/abdd.git
   cd abdd
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Build locally**
   ```bash
   go build -o abdd .
   ```

### Ways to Contribute

- 🐛 **Bug Reports**: Found an issue? Open a GitHub issue with reproduction steps
- 💡 **Feature Requests**: Have an idea? We'd love to hear it!
- 📝 **Documentation**: Help improve our docs and examples
- 🔧 **Code**: Submit pull requests for bug fixes or new features
- 📊 **Examples**: Share real-world test examples

### Reporting Bugs

Please include:
- ABDD version (`abdd --version`)
- Operating system
- Complete error message
- Minimal reproduction case
- Expected vs actual behavior

### Feature Requests

We're particularly interested in:
- Additional validation capabilities  
- CI/CD integration improvements
- Performance optimizations
- Developer experience enhancements

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [GoFakeIt](https://github.com/brianvoe/gofakeit) for providing 310+ fake data generation functions
- [gjson](https://github.com/tidwall/gjson) for powerful JSON path querying
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management

---

**Ready to start testing?** 🚀

```bash
go install github.com/davesavic/abdd@latest
abdd init
abdd run --config abdd.yaml --folders tests
```

Happy testing! 🎯
