// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2021 Datadog, Inc.

package dd_sdk_go_testing

import (
	"fmt"
	"testing"

	"github.com/DataDog/dd-sdk-go-testing/internal/constants"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestStatus(t *testing.T) {
	mt := mocktracer.Start()
	defer mt.Stop()

	t.Run("pass", func(t *testing.T) {
		ctx, finish := StartTest(t)
		defer finish()

		span, _ := tracer.SpanFromContext(ctx)
		span.SetTag("k", "1")
	})

	t.Run("skip", func(t *testing.T) {
		ctx, finish := StartTest(t)
		defer finish()

		span, _ := tracer.SpanFromContext(ctx)
		span.SetTag("k", "2")
		t.Skip("good reason")
	})

	spans := mt.FinishedSpans()
	if len(spans) != 2 {
		t.FailNow()
	}

	const suiteName string = "github.com/DataDog/dd-sdk-go-testing"
	const framework string = "golang.org/pkg/testing"

	s := spans[0]
	assertEqual("test", s.OperationName())
	assertEqual("TestStatus/pass", s.Tag(constants.TestName).(string))
	assertEqual(suiteName, s.Tag(constants.TestSuite).(string))
	assertEqual(fmt.Sprintf("%s.%s", suiteName, "TestStatus/pass"), s.Tag(ext.ResourceName).(string))
	assertEqual(framework, s.Tag(constants.TestFramework).(string))
	assertEqual(constants.TestStatusPass, s.Tag(constants.TestStatus).(string))
	commonEqualCheck(s)
	commonNotEmptyCheck(s)
	fmt.Println(s)

	s = spans[1]
	assertEqual("test", s.OperationName())
	assertEqual("TestStatus/skip", s.Tag(constants.TestName).(string))
	assertEqual(suiteName, s.Tag(constants.TestSuite).(string))
	assertEqual(fmt.Sprintf("%s.%s", suiteName, "TestStatus/skip"), s.Tag(ext.ResourceName).(string))
	assertEqual(framework, s.Tag(constants.TestFramework).(string))
	assertEqual(constants.TestStatusSkip, s.Tag(constants.TestStatus).(string))
	commonEqualCheck(s)
	commonNotEmptyCheck(s)
	fmt.Println(s)
}

func commonEqualCheck(s mocktracer.Span) {
	assertEqual(constants.SpanTypeTest, s.Tag(ext.SpanType).(string))
	assertEqual(constants.SpanTypeTest, s.Tag(constants.SpanKind).(string))
	assertEqual(constants.TestTypeTest, s.Tag(constants.TestType).(string))
	assertEqual(constants.CIAppTestOrigin, s.Tag(constants.Origin).(string))
}

func commonNotEmptyCheck(s mocktracer.Span) {
	assertNotEmpty(s.Tag(constants.GitCommitAuthorDate).(string))
	assertNotEmpty(s.Tag(constants.GitCommitAuthorEmail).(string))
	assertNotEmpty(s.Tag(constants.GitCommitAuthorName).(string))
	assertNotEmpty(s.Tag(constants.GitCommitCommitterDate).(string))
	assertNotEmpty(s.Tag(constants.GitCommitCommitterEmail).(string))
	assertNotEmpty(s.Tag(constants.GitCommitCommitterName).(string))
	assertNotEmpty(s.Tag(constants.GitCommitMessage).(string))
	assertNotEmpty(s.Tag(constants.GitCommitSHA).(string))
	assertNotEmpty(s.Tag(constants.GitRepositoryURL).(string))
	assertNotEmpty(s.Tag(constants.CIWorkspacePath).(string))
	assertNotEmpty(s.Tag(constants.OSArchitecture).(string))
	assertNotEmpty(s.Tag(constants.OSPlatform).(string))
	assertNotEmpty(s.Tag(constants.OSVersion).(string))
}

func assertEqual(expected string, actual string) {
	if expected != actual {
		panic(fmt.Sprintf("Value expected: %s, Actual: %s", expected, actual))
	}
}

func assertNotEmpty(actual string) {
	if actual == "" {
		panic("Value is empty")
	}
}
