
// This file is an automatically generated Java binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package {{.Package}};

import org.ethereum.autonity.*;
import java.util.*;

{{$structs := .Structs}}
{{range $contract := .Contracts}}
{{if not .Library}}public {{end}}class {{.Type}} {
	// ABI is the input ABI used to generate the binding from.
	public final static String ABI = "{{.InputABI}}";
	{{if $contract.FuncSigs}}
		// {{.Type}}FuncSigs maps the 4-byte function signature to its string representation.
		public final static Map<String, String> {{.Type}}FuncSigs;
		static {
			Hashtable<String, String> temp = new Hashtable<String, String>();
			{{range $strsig, $binsig := .FuncSigs}}temp.put("{{$binsig}}", "{{$strsig}}");
			{{end}}
			{{.Type}}FuncSigs = Collections.unmodifiableMap(temp);
		}
	{{end}}
	{{if .InputBin}}
	// BYTECODE is the compiled bytecode used for deploying new contracts.
	public final static String BYTECODE = "0x{{.InputBin}}";

	// deploy deploys a new Ethereum contract, binding an instance of {{.Type}} to it.
	public static {{.Type}} deploy(TransactOpts auth, EthereumClient client{{range .Constructor.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Autonity.newInterfaces({{(len .Constructor.Inputs)}});
		String bytecode = BYTECODE;
		{{if .Libraries}}

		// "link" contract to dependent libraries by deploying them first.
		{{range $pattern, $name := .Libraries}}
		{{capitalise $name}} {{decapitalise $name}}Inst = {{capitalise $name}}.deploy(auth, client);
		bytecode = bytecode.replace("__${{$pattern}}$__", {{decapitalise $name}}Inst.Address.getHex().substring(2));
		{{end}}
		{{end}}
		{{range $index, $element := .Constructor.Inputs}}Interface arg{{$index}} = Autonity.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}
		return new {{.Type}}(Autonity.deployContract(auth, ABI, Autonity.decodeFromHex(bytecode), client, args));
	}

	// Internal constructor used by contract deployment.
	private {{.Type}}(BoundContract deployment) {
		this.Address  = deployment.getAddress();
		this.Deployer = deployment.getDeployer();
		this.Contract = deployment;
	}
	{{end}}

	// Ethereum address where this contract is located at.
	public final Address Address;

	// Ethereum transaction in which this contract was deployed (if known!).
	public final Transaction Deployer;

	// Contract instance bound to a blockchain address.
	private final BoundContract Contract;

	// Creates a new instance of {{.Type}}, bound to a specific deployed contract.
	public {{.Type}}(Address address, EthereumClient client) throws Exception {
		this(Autonity.bindContract(address, ABI, client));
	}

	{{range .Calls}}
	{{if gt (len .Normalized.Outputs) 1}}
	// {{capitalise .Normalized.Name}}Results is the output of a call to {{.Normalized.Name}}.
	public class {{capitalise .Normalized.Name}}Results {
		{{range $index, $item := .Normalized.Outputs}}public {{bindtype .Type $structs}} {{if ne .Name ""}}{{.Name}}{{else}}Return{{$index}}{{end}};
		{{end}}
	}
	{{end}}

	// {{.Normalized.Name}} is a free data retrieval call binding the contract method 0x{{printf "%x" .Original.ID}}.
	//
	// Solidity: {{.Original.String}}
	public {{if gt (len .Normalized.Outputs) 1}}{{capitalise .Normalized.Name}}Results{{else if eq (len .Normalized.Outputs) 0}}void{{else}}{{range .Normalized.Outputs}}{{bindtype .Type $structs}}{{end}}{{end}} {{.Normalized.Name}}(CallOpts opts{{range .Normalized.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Autonity.newInterfaces({{(len .Normalized.Inputs)}});
		{{range $index, $item := .Normalized.Inputs}}Interface arg{{$index}} = Autonity.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}

		Interfaces results = Autonity.newInterfaces({{(len .Normalized.Outputs)}});
		{{range $index, $item := .Normalized.Outputs}}Interface result{{$index}} = Autonity.newInterface(); result{{$index}}.setDefault{{namedtype (bindtype .Type $structs) .Type}}(); results.set({{$index}}, result{{$index}});
		{{end}}

		if (opts == null) {
			opts = Autonity.newCallOpts();
		}
		this.Contract.call(opts, results, "{{.Original.Name}}", args);
		{{if gt (len .Normalized.Outputs) 1}}
			{{capitalise .Normalized.Name}}Results result = new {{capitalise .Normalized.Name}}Results();
			{{range $index, $item := .Normalized.Outputs}}result.{{if ne .Name ""}}{{.Name}}{{else}}Return{{$index}}{{end}} = results.get({{$index}}).get{{namedtype (bindtype .Type $structs) .Type}}();
			{{end}}
			return result;
		{{else}}{{range .Normalized.Outputs}}return results.get(0).get{{namedtype (bindtype .Type $structs) .Type}}();{{end}}
		{{end}}
	}
	{{end}}

	{{range .Transacts}}
	// {{.Normalized.Name}} is a paid mutator transaction binding the contract method 0x{{printf "%x" .Original.ID}}.
	//
	// Solidity: {{.Original.String}}
	public Transaction {{.Normalized.Name}}(TransactOpts opts{{range .Normalized.Inputs}}, {{bindtype .Type $structs}} {{.Name}}{{end}}) throws Exception {
		Interfaces args = Autonity.newInterfaces({{(len .Normalized.Inputs)}});
		{{range $index, $item := .Normalized.Inputs}}Interface arg{{$index}} = Autonity.newInterface();arg{{$index}}.set{{namedtype (bindtype .Type $structs) .Type}}({{.Name}});args.set({{$index}},arg{{$index}});
		{{end}}
		return this.Contract.transact(opts, "{{.Original.Name}}"	, args);
	}
	{{end}}

    {{if .Fallback}}
	// Fallback is a paid mutator transaction binding the contract fallback function.
	//
	// Solidity: {{.Fallback.Original.String}}
	public Transaction Fallback(TransactOpts opts, byte[] calldata) throws Exception {
		return this.Contract.rawTransact(opts, calldata);
	}
    {{end}}

    {{if .Receive}}
	// Receive is a paid mutator transaction binding the contract receive function.
	//
	// Solidity: {{.Receive.Original.String}}
	public Transaction Receive(TransactOpts opts) throws Exception {
		return this.Contract.rawTransact(opts, null);
	}
    {{end}}
}
{{end}}