package ebs

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/packer/packer"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"access_key":    "foo",
		"secret_key":    "bar",
		"source_ami":    "foo",
		"instance_type": "foo",
		"region":        "us-east-1",
		"ssh_username":  "root",
		"ami_name":      "foo",
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{}
	raw = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilder_Prepare_BadType(t *testing.T) {
	b := &Builder{}
	c := map[string]interface{}{
		"access_key": []string{},
	}

	warnings, err := b.Prepare(c)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("prepare should fail")
	}
}

func TestBuilderPrepare_AMIName(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test good
	config["ami_name"] = "foo"
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test bad
	config["ami_name"] = "foo {{"
	b = Builder{}
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}

	// Test bad
	delete(config, "ami_name")
	b = Builder{}
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_InvalidKey(t *testing.T) {
	var b Builder
	config := testConfig()

	// Add a random key
	config["i_should_not_be_valid"] = true
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_InvalidShutdownBehavior(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test good
	config["shutdown_behavior"] = "terminate"
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test good
	config["shutdown_behavior"] = "stop"
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test bad
	config["shutdown_behavior"] = "foobar"
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func createTestStateBagStepStopEbsInstance() multistep.StateBag {
	// Make a faked UI, instance, and ec2 conection
	var out, err bytes.Buffer
	var ui packer.Ui = &packer.BasicUi{
		Writer:      &out,
		ErrorWriter: &err,
	}
	FakeInstance := &ec2.Instance{
		InstanceId: aws.String("instance-id"),
	}

	type FakeEC2Conn struct {
		*string
	}

	func (c FakeEC2Conn) StopInstances(input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
		return nil, nil
	}

	func (c FakeEC2Conn) StopInstances()  {
		return c.FakeStopInstances()
	}

	var mockEC2Conn FakeEC2Conn

	// Set up state bag for test using generated state.
	state := new(multistep.BasicStateBag)
	state.Put("ec2", mockEC2Conn)
	state.Put("ui", ui)
	state.Put("instance", FakeInstance)

	// These are teh states grabbed by step_stop_ebs_instance
	// ec2conn := state.Get("ec2").(*ec2.EC2)
	// instance := state.Get("instance").(*ec2.Instance)
	// ui := state.Get("ui").(packer.Ui)

}
