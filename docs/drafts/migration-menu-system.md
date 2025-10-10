# Menu System Refactoring Migration Plan

## Problem Statement

The current menu system mixes business logic with presentation concerns, creating several architectural issues:

1. **Mixed Concerns**: Web handlers contain presentation logic (`menu.AddNewItem()`, `menu.AddListItem()`) that should belong in templates
2. **Configuration Smell**: Button styles are defined via environment variables (`CLIO_BUTTON_STYLE_*`) instead of proper CSS classes
3. **Scattered Logic**: Menu construction logic is duplicated across 6 handlers (21 occurrences total)
4. **CSS Anti-pattern**: Raw CSS classes passed through configuration instead of semantic class names

Current problematic pattern:
```go
// Handler mixing business + presentation logic
menu := page.NewMenu(ssgPath)
menu.AddNewItem(&Content{})     // "I need a New button"
menu.AddListItem(&content, "Back") // "I need a Back button"
```

## Proposed Solution

Migrate to a **declarative template-based menu system** with:

1. **Template Function Injection**: Extend the generic `hm.TemplateManager` to accept custom functions
2. **Declarative Menus**: Move menu definitions to templates where they belong  
3. **Semantic CSS Classes**: Replace raw CSS configuration with proper design system classes
4. **Clean Separation**: Remove all presentation logic from handlers

## Architecture Changes

### 1. Template Manager Enhancement

Extend `hm.TemplateManager` with function injection capability:

```go
// In hm package - remains generic
type TemplateManager struct {
    baseFuncMap   template.FuncMap  // Core functions (Truncate, Add, etc.)
    customFuncMap template.FuncMap  // App-injected functions
}

func (tm *TemplateManager) RegisterFunctions(funcs template.FuncMap) {
    // Merge custom functions for template compilation
}
```

### 2. SSG-Specific Function Registration

Register domain-specific template functions in the SSG module:

```go
// In internal/web/ssg - app-specific
ssgFunctions := template.FuncMap{
    "newPath":  func(entityType string) string { /* /ssg/new-{entity} */ },
    "listPath": func(entityType string) string { /* /ssg/{entity}s */ },
    "editPath": func(entityType, id string) string { /* /ssg/edit-{entity}?id={id} */ },
}
tm.RegisterFunctions(ssgFunctions)
```

### 3. Template-Based Menu Declarations

Replace programmatic menu construction with declarative templates:

```html
<!-- BEFORE: Handler-driven -->
{{ range .Menu.Items }}
  <a href="{{ .Path }}" class="{{ .Style }}">{{ .Text }}</a>
{{ end }}

<!-- AFTER: Template-driven -->
{{ define "submenu" }}
<div class="mx-auto p-4">
  <div class="flex space-x-4 justify-center">
    <a href="{{ newPath "content" }}" class="btn btn-primary">New Content</a>
    <a href="{{ listPath "content" }}" class="btn btn-secondary">Back to List</a>
  </div>
</div>
{{ end }}
```

### 4. CSS Design System

Replace environment-based styles with semantic CSS classes:

```css
/* Semantic button system using Tailwind @apply */
.btn {
  @apply px-4 py-2 rounded font-medium transition-colors;
}
.btn-primary { @apply bg-blue-600 text-white hover:bg-blue-700; }
.btn-secondary { @apply bg-gray-600 text-white hover:bg-gray-700; }
.btn-danger { @apply bg-red-600 text-white hover:bg-red-700; }
.btn-success { @apply bg-green-600 text-white hover:green-700; }
```

## Migration Strategy

### Phase 1: Foundation Setup
1. **Extend Template Manager**: Add function injection capability to `hm.TemplateManager`
2. **CSS Classes**: Define semantic button classes in CSS
3. **Template Functions**: Register SSG-specific path generation functions

### Phase 2: Incremental Entity Migration
Migrate one entity at a time to avoid breaking changes:

#### Per Entity Steps:
1. **Update Template**: Replace `{{ template "menu" . }}` with declarative menu HTML
2. **Test Template**: Verify correct buttons appear with proper styling
3. **Clean Handler**: Remove `menu.AddXXXItem()` calls from corresponding handler
4. **Verify Functionality**: Ensure all menu actions work correctly

#### Migration Order:
1. **Content** (most used, best test case)
2. **Image** 
3. **Section**
4. **Layout**
5. **Tag**
6. **Param**

### Phase 3: Cleanup
1. **Remove Environment Variables**: Delete `CLIO_BUTTON_STYLE_*` from makefile
2. **Clean Menu Types**: Remove unused menu-related code from handlers
3. **Update Documentation**: Document new template function system

## Files Affected

### Core Changes:
- `hm/template.go`: Function injection capability
- `internal/web/ssg/webhandler.go`: Function registration
- `assets/static/css/main.css`: Semantic button classes

### Handler Cleanup (6 files):
- `internal/web/ssg/webhandlercontent.go`
- `internal/web/ssg/webhandlerimage.go`  
- `internal/web/ssg/webhandlerlayout.go`
- `internal/web/ssg/webhandlerparam.go`
- `internal/web/ssg/webhandlersection.go`
- `internal/web/ssg/webhandlertag.go`

### Template Updates:
- All `assets/template/handler/ssg/*.tmpl` files with menu definitions

## Implementation Notes

The migration is designed to be incremental, one entity at a time to avoid breaking everything at once. Template manager changes are purely additive, so existing functionality continues working during the transition.

We are going to start with Content as it's the most used entity. If that works cleanly, the rest should be mechanical.
