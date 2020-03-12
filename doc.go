// Package reply package is essentially a source code conversion
// of the ruby library https://github.com/discourse/email_reply_trimmer.
// The core logic is a almost line by line conversion.
//
// This package has a dependency on excellent regex library github.com/dlclark/regexp2.
// The reason for not using the standard regex library was due to the fact that
// the regex package from the stdlib is not compatible with the library from the Ruby stdlib.
//
// All the tests were taken from the email_reply_trimmer library.
//
// Note:
// This code is not idiomatic go code, as, it was mostly adapted from the ruby code,
// however, the public APIs were kept simple as possible and does not expose any internal.
package reply
