// Модуль подключает линтеры к приложению.
/*
# Пример запуска линтера
	go run cmd/staticlint/main.go -all ./...
Линтеры можно запускать отдельно
Подробнее:
	go run cmd/staticlint/main.go help

Подключенные линтеры:


# staticcheck : the advanced Go linter.

	https://staticcheck.dev/

# asmdecl : report mismatches between assembly files and Go declarations


# assign : check for useless assignments

This checker reports assignments of the form x = x or a[i] = a[i].
These are almost always useless, and even when they aren't they are
usually a mistake.


# atomic : check for common mistakes using the sync/atomic package

The atomic checker looks for assignment statements of the form:

        x = atomic.AddUint64(&x, 1)

which are not atomic.


# atomicalign : check for non-64-bits-aligned arguments to sync/atomic functions


# bools : check for common mistakes involving boolean operators


# buildssa : build SSA-form IR for later passes


# buildtag : check //go:build and // +build directives


# cgocall : detect some violations of the cgo pointer passing rules

Check for invalid cgo pointer passing.
This looks for code that uses cgo to call C code passing values
whose types are almost always invalid according to the cgo pointer
sharing rules.
Specifically, it warns about attempts to pass a Go chan, map, func,
or slice to C, either directly, or via a pointer, array, or struct.


# composites : check for unkeyed composite literals

This analyzer reports a diagnostic for composite literals of struct
types imported from another package that do not use the field-keyed
syntax. Such literals are fragile because the addition of a new field
(even if unexported) to the struct will cause compilation to fail.

As an example,

        err = &net.DNSConfigError{err}

should be replaced by:

        err = &net.DNSConfigError{Err: err}



# copylocks : check for locks erroneously passed by value

Inadvertently copying a value containing a lock, such as sync.Mutex or
sync.WaitGroup, may cause both copies to malfunction. Generally such
values should be referred to through a pointer.


# ctrlflow : build a control-flow graph


# deepequalerrors : check for calls of reflect.DeepEqual on error values

The deepequalerrors checker looks for calls of the form:

    reflect.DeepEqual(err1, err2)

where err1 and err2 are errors. Using reflect.DeepEqual to compare
errors is discouraged.


# defer : report common mistakes in defer statements

The defer analyzer reports a diagnostic when a defer statement would
result in a non-deferred call to time.Since, as experience has shown
that this is nearly always a mistake.

For example:

        start := time.Now()
        ...
        defer recordLatency(time.Since(start)) // error: call to time.Since is not deferred

The correct code is:

        defer func() { recordLatency(time.Since(start)) }()


# directive : check Go toolchain directives such as //go:debug

# errorsas : report passing non-pointer or non-error values to errors.As

# fieldalignment : find structs that would use less memory if their fields were sorted

This analyzer find structs that can be rearranged to use less memory, and provides
a suggested edit with the most compact order.

Note that there are two different diagnostics reported. One checks struct size,
and the other reports "pointer bytes" used. Pointer bytes is how many bytes of the
object that the garbage collector has to potentially scan for pointers, for example:

        struct { uint32; string }

have 16 pointer bytes because the garbage collector has to scan up through the string's
inner pointer.

        struct { string; *uint32 }

has 24 pointer bytes because it has to scan further through the *uint32.

        struct { string; uint32 }

has 8 because it can stop immediately after the string pointer.

Be aware that the most compact order is not always the most efficient.
In rare cases it may cause two variables each updated by its own goroutine
to occupy the same CPU cache line, inducing a form of memory contention
known as "false sharing" that slows down both goroutines.



# findcall : find calls to a particular function

The findcall analysis reports calls to functions or methods
of a particular name.


# framepointer : report assembly that clobbers the frame pointer before saving it


# httpresponse : check for mistakes using HTTP responses

A common mistake when using the net/http package is to defer a function
call to close the http.Response Body before checking the error that
determines whether the response is valid:

        resp, err := http.Head(url)
        defer resp.Body.Close()
        if err != nil {
                log.Fatal(err)
        }
        // (defer statement belongs here)

This checker helps uncover latent nil dereference bugs by reporting a
diagnostic for such mistakes.


# ifaceassert : detect impossible interface-to-interface type assertions

This checker flags type assertions v.(T) and corresponding type-switch cases
in which the static type V of v is an interface that cannot possibly implement
the target interface T. This occurs when V and T contain methods with the same
name but different signatures. Example:

        var v interface {
                Read()
        }
        _ = v.(io.Reader)

The Read method in v has a different signature than the Read method in
io.Reader, so this assertion cannot succeed.


# inspect : optimize AST traversal for later passes


# loopclosure : check references to loop variables from within nested functions

This analyzer reports places where a function literal references the
iteration variable of an enclosing loop, and the loop calls the function
in such a way (e.g. with go or defer) that it may outlive the loop
iteration and possibly observe the wrong value of the variable.

In this example, all the deferred functions run after the loop has
completed, so all observe the final value of v.

        for _, v := range list {
            defer func() {
                use(v) // incorrect
            }()
        }

One fix is to create a new variable for each iteration of the loop:

        for _, v := range list {
            v := v // new var per iteration
            defer func() {
                use(v) // ok
            }()
        }

The next example uses a go statement and has a similar problem.
In addition, it has a data race because the loop updates v
concurrent with the goroutines accessing it.

        for _, v := range elem {
            go func() {
                use(v)  // incorrect, and a data race
            }()
        }

A fix is the same as before. The checker also reports problems
in goroutines started by golang.org/x/sync/errgroup.Group.
A hard-to-spot variant of this form is common in parallel tests:

        func Test(t *testing.T) {
            for _, test := range tests {
                t.Run(test.name, func(t *testing.T) {
                    t.Parallel()
                    use(test) // incorrect, and a data race
                })
            }
        }

The t.Parallel() call causes the handlers_rest of the function to execute
concurrent with the loop.

The analyzer reports references only in the last statement,
as it is not deep enough to understand the effects of subsequent
statements that might render the reference benign.
("Last statement" is defined recursively in compound
statements such as if, switch, and select.)

See: https://golang.org/doc/go_faq.html#closures_and_goroutines


# lostcancel : check cancel func returned by context.WithCancel is called

The cancellation function returned by context.WithCancel, WithTimeout,
and WithDeadline must be called or the new context will remain live
until its parent context is cancelled.
(The background context is never cancelled.)


# nilfunc : check for useless comparisons between functions and nil

A useless comparison is one like f == nil as opposed to f() == nil.


# nilness : Annotates return values that will never be nil (typed or untyped)


# pkgfact : gather name/value pairs from constant declarations


# printf : check consistency of Printf format strings and arguments

The check applies to calls of the formatting functions such as
[fmt.Printf] and [fmt.Sprintf], as well as any detected wrappers of
those functions.

In this example, the %d format operator requires an integer operand:

        fmt.Printf("%d", "hello") // fmt.Printf format %d has arg "hello" of wrong type string

See the documentation of the fmt package for the complete set of
format operators and their operand types.

To enable printf checking on a function that is not found by this
analyzer's heuristics (for example, because control is obscured by
dynamic method calls), insert a bogus call:

        func MyPrintf(format string, args ...any) {
                if false {
                        _ = fmt.Sprintf(format, args...) // enable printf checker
                }
                ...
        }

The -funcs flag specifies a comma-separated list of names of additional
known formatting functions or methods. If the name contains a period,
it must denote a specific function using one of the following forms:

        dir/pkg.Function
        dir/pkg.Type.Method
        (*dir/pkg.Type).Method

Otherwise the name is interpreted as a case-insensitive unqualified
identifier such as "errorf". Either way, if a listed name ends in f, the
function is assumed to be Printf-like, taking a format string before the
argument list. Otherwise it is assumed to be Print-like, taking a list
of arguments with no format string.


# reflectvaluecompare : check for comparing reflect.Value values with == or reflect.DeepEqual

The reflectvaluecompare checker looks for expressions of the form:

        v1 == v2
        v1 != v2
        reflect.DeepEqual(v1, v2)

where v1 or v2 are reflect.Values. Comparing reflect.Values directly
is almost certainly not correct, as it compares the reflect package's
internal representation, not the underlying value.
Likely what is intended is:

        v1.Interface() == v2.Interface()
        v1.Interface() != v2.Interface()
        reflect.DeepEqual(v1.Interface(), v2.Interface())

# shift : check for shifts that equal or exceed the width of the integer


# sigchanyzer : check for unbuffered channel of os.Signal

This checker reports call expression of the form

        signal.Notify(c <-chan os.Signal, sig ...os.Signal),

where c is an unbuffered channel, which can be at risk of missing the signal.


# slog : check for invalid structured logging calls

The slog checker looks for calls to functions from the log/slog
package that take alternating key-value pairs. It reports calls
where an argument in a key position is neither a string nor a
slog.Attr, and where a final key is missing its value.
For example,it would report

        slog.Warn("message", 11, "k") // slog.Warn arg "11" should be a string or a slog.Attr

and

        slog.Info("message", "k1", v1, "k2") // call to slog.Info missing a final value


# sortslice : check the argument type of sort.Slice

# stdmethods : check signature of methods of well-known interfaces

Sometimes a type may be intended to satisfy an interface but may fail to
do so because of a mistake in its method signature.
For example, the result of this WriteTo method should be (int64, error),
not error, to satisfy io.WriterTo:

        type myWriterTo struct{...}
        func (myWriterTo) WriteTo(w io.Writer) error { ... }

This check ensures that each method whose name matches one of several
well-known interface methods from the standard library has the correct
signature for that interface.

Checked method names include:

        Format GobEncode GobDecode MarshalJSON MarshalXML
        Peek ReadByte ReadFrom ReadRune Scan Seek
        UnmarshalJSON UnreadByte UnreadRune WriteByte
        WriteTo


# stringintconv : check for string(int) conversions

This checker flags conversions of the form string(x) where x is an integer
(but not byte or rune) type. Such conversions are discouraged because they
return the UTF-8 representation of the Unicode code point x, and not a decimal
string representation of x as one might expect. Furthermore, if x denotes an
invalid code point, the conversion cannot be statically rejected.

For conversions that intend on using the code point, consider replacing them
with string(rune(x)). Otherwise, strconv.Itoa and its equivalents return the
string representation of the value in the desired base.


# structtag : check that struct field tags conform to reflect.StructTag.Get

Also report certain struct tags (json, xml) used with unexported fields.


# testinggoroutine : report calls to (*testing.T).Fatal from goroutines started by a test.

Functions that abruptly terminate a test, such as the Fatal, Fatalf, FailNow, and
Skip{,f,Now} methods of *testing.T, must be called from the test goroutine itself.
This checker detects calls to these functions that occur within a goroutine
started by the test. For example:

        func TestFoo(t *testing.T) {
            go func() {
                t.Fatal("oops") // error: (*T).Fatal called from non-test goroutine
            }()
        }


# tests : check for common mistaken usages of tests and examples

The tests checker walks Test, Benchmark, Fuzzing and Example functions checking
malformed names, wrong signatures and examples documenting non-existent
identifiers.

Please see the documentation for package testing in golang.org/pkg/testing
for the conventions that are enforced for Tests, Benchmarks, and Examples.


# timeformat : check for calls of (time.Time).Format or time.Parse with 2006-02-01

The timeformat checker looks for time formats with the 2006-02-01 (yyyy-dd-mm)
format. Internationally, "yyyy-dd-mm" does not occur in common calendar date
standards, and so it is more likely that 2006-01-02 (yyyy-mm-dd) was intended.


# unmarshal : report passing non-pointer or non-interface values to unmarshal

The unmarshal analysis reports calls to functions such as json.Unmarshal
in which the argument type is not a pointer or an interface.


# unreachable : check for unreachable code

The unreachable analyzer finds statements that execution can never reach
because they are preceded by an return statement, a call to panic, an
infinite loop, or similar constructs.


# unsafeptr : check for invalid conversions of uintptr to unsafe.Pointer

The unsafeptr analyzer reports likely incorrect uses of unsafe.Pointer
to convert integers to pointers. A conversion from uintptr to
unsafe.Pointer is invalid if it implies that there is a uintptr-typed
word in memory that holds a pointer value, because that word will be
invisible to stack copying and to the garbage collector.


# unusedresult : check for unused results of calls to some functions

Some functions like fmt.Errorf return a result and have no side
effects, so it is always a mistake to discard the result. Other
functions may return an error that must not be ignored, or a cleanup
operation that must be called. This analyzer reports calls to
functions like these when the result of the call is ignored.

The set of functions may be controlled using flags.


# unusedwrite : checks for unused writes

The analyzer reports instances of writes to struct fields and
arrays that are never read. Specifically, when a struct object
or an array is copied, its elements are copied implicitly by
the compiler, and any element write to this copy does nothing
with the original object.

For example:

        type T struct { x int }

        func f(input []T) {
                for i, v := range input {  // v is a copy
                        v.x = i  // unused write to field x
                }
        }

Another example is about non-pointer receiver:

        type T struct { x int }

        func (t T) f() {  // t is a copy
                t.x = i  // unused write to field x
        }


# usesgenerics : detect whether a package uses generics features

The usesgenerics analysis reports whether a package directly or transitively
uses certain features associated with generic programming in Go.


# osExit : check for os.Exit in main


# errcheck : check for unchecked errors


# ruleguard : The most opinionated Go source code linter.

*/
package main
