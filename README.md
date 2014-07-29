# structfield 

`structfield` provides an API to perform processing on the fields of a `struct` which will be JSON-marshalled.
This is useful for e.g:

 * change a field name based on an external decision,
 * discard unwanted fields
 * and in general, apply logic which isn't possible or is conflicting with json tags.
