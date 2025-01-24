// Code generated for internal testing purposes only - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"math/big"
	"strings"
	"errors"
	"fmt"

	ethereum "github.com/autonity/autonity"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

{{$structs := .Structs}}
{{range $structs}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{$field.Name}} {{$field.Type}}{{end}}
	}
{{end}}

{{range $contract := .Contracts}}
	// {{.Type}}MetaData contains all meta data concerning the {{.Type}} contract.
	var {{.Type}}MetaData = &bind.MetaData{
		ABI: "{{.InputABI}}",
		{{if $contract.FuncSigs -}}
		Sigs: map[string]string{
			{{range $strsig, $binsig := .FuncSigs}}"{{$binsig}}": "{{$strsig}}",
			{{end}}
		},
		{{end -}}
		{{if .InputBin -}}
		Bin: "0x{{.InputBin}}",
		{{end}}
	}
	// {{.Type}}ABI is the input ABI used to generate the binding from.
	// Deprecated: Use {{.Type}}MetaData.ABI instead.
	var {{.Type}}ABI = {{.Type}}MetaData.ABI

	{{if $contract.FuncSigs}}
		// Deprecated: Use {{.Type}}MetaData.Sigs instead.
		// {{.Type}}FuncSigs maps the 4-byte function signature to its string representation.
		var {{.Type}}FuncSigs = {{.Type}}MetaData.Sigs
	{{end}}

	{{if .InputBin}}
		// {{.Type}}Bin is the compiled bytecode used for deploying new contracts.
		// Deprecated: Use {{.Type}}MetaData.Bin instead.
		var {{.Type}}Bin = {{.Type}}MetaData.Bin

		// Deploy{{.Type}} deploys a new Ethereum contract, binding an instance of {{.Type}} to it.
		func (r *Runner) Deploy{{.Type}}(opts *runOptions {{range .Constructor.Inputs}}, {{.Name}} {{bindtype .Type $structs}}{{end}}) (common.Address, uint64, *{{.Type}}, error) {
		  parsed, err := {{.Type}}MetaData.GetAbi()
		  if err != nil {
		    return common.Address{}, 0, nil, err
		  }
		  if parsed == nil {
			return common.Address{}, 0, nil, errors.New("GetABI returned nil")
		  }
		  {{range $pattern, $name := .Libraries}}
            /* LIBRARY DEPLOYMENT NOT SUPPORTED	*/
			{{decapitalise $name}}Addr, _, _, _ := Deploy{{capitalise $name}}(auth, backend)
			{{$contract.Type}}Bin = strings.Replace({{$contract.Type}}Bin, "__${{$pattern}}$__", {{decapitalise $name}}Addr.String()[2:], -1)
		  {{end}}
		  address, gasConsumed, c, data, err := r.deployContract(opts, parsed, common.FromHex({{.Type}}Bin) {{range .Constructor.Inputs}}, {{.Name}}{{end}})
		  if err != nil {
		    return common.Address{}, 0, nil, (&{{.Type}}{contract: c}).DecodeError(data, err)
		  }
		  return address, gasConsumed, &{{.Type}}{contract: c}, nil
		}
	{{end}}

	// {{.Type}} is an auto generated Go binding around an Ethereum contract.
	type {{.Type}} struct {
		*contract
	}

	{{range .Calls}}
		// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}) {{.Normalized.Name}}(opts *runOptions {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} uint64, error) {
			data, consumed, err := _{{$contract.Type}}.call(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			{{if .Structured}}
			outstruct := new(struct{ {{range .Normalized.Outputs}} {{.Name}} {{bindtype .Type $structs}}; {{end}} })
			if err != nil {
				return *outstruct, consumed, _{{$contract.Type}}.DecodeError(data, err)
			}
			out, err := _{{$contract.Type}}.abi.Unpack("{{.Original.Name}}", data)
			if err != nil {
                return *outstruct, consumed, err
            }
			{{range $i, $t := .Normalized.Outputs}}
			outstruct.{{.Name}} = *abi.ConvertType(out[{{$i}}], new({{bindtype .Type $structs}})).(*{{bindtype .Type $structs}}){{end}}
			return *outstruct, consumed, err
			{{else}}
			if err != nil {
				return {{range $i, $_ := .Normalized.Outputs}}*new({{bindtype .Type $structs}}), {{end}} consumed, _{{$contract.Type}}.DecodeError(data, err)
			}
			out, err := _{{$contract.Type}}.abi.Unpack("{{.Original.Name}}", data)
			if err != nil {
                return {{range $i, $_ := .Normalized.Outputs}}*new({{bindtype .Type $structs}}), {{end}} consumed, err
            }
			{{range $i, $t := .Normalized.Outputs}}
			out{{$i}} := *abi.ConvertType(out[{{$i}}], new({{bindtype .Type $structs}})).(*{{bindtype .Type $structs}}){{end}}
			return {{range $i, $t := .Normalized.Outputs}}out{{$i}}, {{end}} consumed, err
			{{end}}
		}
	{{end}}

	{{range .Transacts}}
		// {{.Normalized.Name}} is a free data retrieval call for a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
		// Similar to eth_call from rpc calls or function.call from truffle, it reverts the state after the call and returns the output. The output is extracted
		// the same way as done above for view only functions.
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}) Call{{.Normalized.Name}}(r *Runner, opts *runOptions {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) ({{if .Structured}}struct{ {{range .Normalized.Outputs}}{{.Name}} {{bindtype .Type $structs}};{{end}} },{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}},{{end}}{{end}} uint64, error) {
			snap := r.snapshot()
			{{if not .Normalized.Outputs}}
			data, consumed, err := _{{$contract.Type}}.call(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			r.revertSnapshot(snap)
			return consumed, _{{$contract.Type}}.DecodeError(data, err)
			{{else}}
			data, consumed, err := _{{$contract.Type}}.call(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			r.revertSnapshot(snap)
			{{if .Structured}}
			outstruct := new(struct{ {{range .Normalized.Outputs}} {{.Name}} {{bindtype .Type $structs}}; {{end}} })
			if err != nil {
				return *outstruct, consumed, _{{$contract.Type}}.DecodeError(data, err)
			}
			out, err := _{{$contract.Type}}.abi.Unpack("{{.Original.Name}}", data)
			if err != nil {
                return *outstruct, consumed, err
            }
			{{range $i, $t := .Normalized.Outputs}}
			outstruct.{{.Name}} = *abi.ConvertType(out[{{$i}}], new({{bindtype .Type $structs}})).(*{{bindtype .Type $structs}}){{end}}
			return *outstruct, consumed, err
			{{else}}
			if err != nil {
				return {{range $i, $_ := .Normalized.Outputs}}*new({{bindtype .Type $structs}}), {{end}} consumed, _{{$contract.Type}}.DecodeError(data, err)
			}
			out, err := _{{$contract.Type}}.abi.Unpack("{{.Original.Name}}", data)
			if err != nil {
                return {{range $i, $_ := .Normalized.Outputs}}*new({{bindtype .Type $structs}}), {{end}} consumed, err
            }
			{{range $i, $t := .Normalized.Outputs}}
			out{{$i}} := *abi.ConvertType(out[{{$i}}], new({{bindtype .Type $structs}})).(*{{bindtype .Type $structs}}){{end}}
			return {{range $i, $t := .Normalized.Outputs}}out{{$i}}, {{end}} consumed, err
			{{end}}
			{{end}}
		}
	{{end}}

	{{range .Transacts}}
		// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
		//
		// Solidity: {{.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}) {{.Normalized.Name}}(opts *runOptions {{range .Normalized.Inputs}}, {{.Name}} {{bindtype .Type $structs}} {{end}}) (uint64, error) {
			data, consumed, err := _{{$contract.Type}}.call(opts, "{{.Original.Name}}" {{range .Normalized.Inputs}}, {{.Name}}{{end}})
			return consumed, _{{$contract.Type}}.DecodeError(data, err)
		}
	{{end}}

	{{if .Fallback}}
		// Fallback is a paid mutator transaction binding the contract fallback function.
		// WARNING! UNTESTED
		// Solidity: {{.Fallback.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}) Fallback(opts *runOptions, calldata []byte) (uint64, error) {
			out, consumed, err := _{{$contract.Type}}.call(opts, "", calldata)
			return consumed, _{{$contract.Type}}.DecodeError(out, err)
		}
	{{end}}

	{{if .Receive}}
		// Receive is a paid mutator transaction binding the contract receive function.
		// WARNING! UNTESTED
		// Solidity: {{.Receive.Original.String}}
		func (_{{$contract.Type}} *{{$contract.Type}}) Receive(opts *runOptions) (uint64, error) {
			out, consumed, err := _{{$contract.Type}}.call(opts, "")
			return consumed, _{{$contract.Type}}.DecodeError(out, err)
		}
	{{end}}

    {{range .Errors}}
        // {{.Normalized.Name}} is a contract revert error corresponding to contract error 0x{{.Selector}}.
    	type {{$contract.Type}}{{.Normalized.Name}}Error struct {
            {{range .Normalized.Inputs}}
            {{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}; {{end}}
        }

        func (e {{$contract.Type}}{{.Normalized.Name}}Error) Error() string {
            return "execution reverted: {{.Normalized.Name}}"
        }
    {{end}}

    func (_{{$contract.Type}} *{{$contract.Type}}) DecodeError(data []byte, err error) error {
        if err == nil {
            return nil
        }
        {{if not .Errors}}
        reason, _ := abi.UnpackRevert(data)
        return fmt.Errorf("%w: %s", err, reason)
        {{else}}
        if len(data) < 4 {
            return err
        }
        selector := data[:4]
        switch common.Bytes2Hex(selector) {
        {{range .Errors}}
        case "{{.Selector}}":
            var e {{$contract.Type}}{{.Normalized.Name}}Error
            err := _{{$contract.Type}}.abi.UnpackIntoInterface(&e, "{{.Original.Name}}", data)
            if err != nil {
                return err
            }
            return e
        {{end}}
        default:
            reason, _ := abi.UnpackRevert(data)
            return fmt.Errorf("%w: %s", err, reason)
        }
        {{end}}
    }
{{end}}