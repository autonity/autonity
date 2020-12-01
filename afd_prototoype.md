## Notation

pi → The participant who's message is being checked
⊥ → No message sent (this is new, to solve problems with the rules)

## Interpreting the rules

Looking at page 11 of the document it is explained that if we are seeing a
message from a pariticipant this implies a certain state of non observable
variables at pi, and from that we can reason what messages must have caused
those state variables to be set at pi, if we cannot see those messages then we
can consider pi to have misbehaved. It also states that the contrary is not
true seeing the messgaes that could precipitate a state change does not
imply the sending of a message, becasue we cannot be sure that they were
recived in a timely manner by other participants.

What does this mean for the rules? Basically if we see the right side, then we
expect to see the corresponding left side if we cannot see the left side then we
have bad behaviour.

It should be noted that it is assumed that when applying the rules all relevant
messages have been received. In practice nodes can never be sure that they will
have received all the messgaes, but nodes can wait some period to maximise the
chance that they do receive all messages.

## Broken rules

The rules PN, PO, PVN and PVO are broken because they all contain terms of the
form "(...|pi)\*". The meaning of the asterisk is that the content of the term
expands to "fill" all the rounds between the rounds referenced from other
terms. The problem with that is the fact that nodes can move to the next round
without sending any messages (see line 55 tendermint pseudocode), wheras these
"(...|pi)\*" terms imply that pi sends a message in all rounds that the term
has expanded to "fill".

Additionally the rules that contain `V ⊕ nil` on the right hand side.

So before we start we must fix the rules.

### PN

PN:  `(M r'<r, PC|pi )* ⇒  M r,P|pi`

#### PN1
Old PN1: `nil ⇒  V` 

The old PN1 says that in all rounds before the round where pi proposes V as a
new value we should see a precommit for nil from pi.

Fixed PN1: `nil v ⊥ ⇒  V` 

The fixed PN1 says that in all rounds before the round where pi proposes V as a
new value we should see either no precommit from pi or a precommit for nil from
pi.

### PO
PO:	 `M r'<r,PV ≺ (M r'≤r"<r,PC|pi)* ⇒  M r,P|pi`

#### PO1
Old PO1: `#(M r' ,PV|V) ≥ 2f + 1 ≺  nil v  V ⇒  V` 

The old PO1 says that if we see pi propose an old V then we should see a
previous quorum of prevotes for V and in all intervening rounds till the
proposal, pi should precommit for nil or V. 

I'm not sure we can provide a good fix for PO1 because what we want it to say is:

That from the most recently observed quorum of prevotes for V we should only
witness precommits for V precommits for nil or no precommit from pi, in the
intervening rounds till pi proposes V. It's difficult to express this in a rule
since we have no notation to mark the most recent of something.

Without significantly chaning the notation for the rules I think the best we
can do for PO1 is this.

Fixed PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ⇒  V` 

The fixed PO1 says if we see a proposal from pi for a new value V then we
expect to see a quorum of precommits for V in the past.

### PVN
PVN: `M r'<r,PC|pi ≺  (M r'<r"<r,PC|pi)* ≺  M r,P ⇒  M r,PV|pi`

#### PVN2
Old PVN2:  `r'=0 ∧ (nil ≺ nil ≺ V |Valid(V)) ⇒ V ⊕ nil`

This rule just doesn't make sense to me at all? There are two possible messgaes
we could observe from the right side one is nil (apparently from a timeout) but
seeing a nil message doesn't allow you to infer anything about the state of the
node. The rule also does not make sense in and of itself since it is called
"prevote new" only prevotes for V can be a prevote for a new value a prevote
for nil is not a prevote for anything.

It actually looks like the researchers forgot which way round their rules work
because the description of this rule talks about what the node would send as a
prevote based on the prior seen messgaes. Instead we need to say "given this
message that implies a certain state, what messages would we expect to see that
would cause that state to come about".

The nil prevote is apparently in case of a timeout, but if a node prevotes nil,
it you cannot expect anything prior to that. So if the remove the nil does the
rule make sense?

Take2 PVN2:  `r'=0 ∧ (nil ≺ nil ≺ V |Valid(V)) ⇒ V`

This says if we see a prevote for a new value v from pi then in all rounds
prior pi sent a precommit for nil, still not right.

Now we need to fix the use of "(...|pi)\*".

Take3 PVN2:  `r'=0 ∧ ( nil v ⊥  ≺  nil v ⊥  ≺  V|Valid(V) ) ⇒ V`

This says if we see a prevote for a new value V then in all previous rounds we
expect to see no precommit or a precommit for nil from pi.

This still looks wrong since it is valid for pi to have a locked value of V and
prevote for a new value (line 23 tendermint).

### PVO

I see the same issues for PVO as for PVN. The use of `V ⊕ nil` doesn't make
sense.

#### PVN1

# How can rules be violated?

we use pi to define the participant under scrutiny. And we understand the
meaning of * when applied to an expression in brackets to mean that that
expression expands to fill the range of rounds defined in the brackets but that
that range of rounds may be empty.

Looking at page 11 of the document it is explained that if we are seeing a
message from a pariticipant this implies a certain state of non observable
variables at pi, and from that we can reason what messages must have caused
those state variables to be set at pi, if we cannot see those messages then we
can consider pi to have misbehaved. It also states that the contrary is not
true seeing the messgaes that could precipitate a state change does not
imply the sending of a message, becasue we cannot be sure that they were
recived in a timely manner by other participants.

What does this mean for the rules? Basically if we see the right side, then we
expect to see the corresponding left side if we cannot see the left side then we
have bad behaviour.

It should be noted that it is assumed that when applying the rules all relevant
messages have been received. In practice nodes can never be sure that they will
have received all the messgaes, but nodes can wait some period to maximise the
chance that they do receive all messages.

As an example lets take PN1: How can rule PN1 be violated?

So PN1 is violated when a proposal for a new value is sent and we cannot see a
precommit for nil in all previous rounds.

(I realise now that there is a mistake in PN1 because it assumes that seeing a
proposal for a new value implies that pi sent a precommit for nil in all
previous rounds, when in fact the rule at line 55 shows that pi may not send
any precommits in previous rounds.)

So lets fix rule PN1. If we see a proposal for a new value what do we expect to
see in previous rounds? We expect to see precommits for nil or no precommits at
all.

So a deviation from this would be seeing a precommit for any value from pi.

This is what I will call a provable violation, since we instantly know that pi
has misbehaved. 

There are however cases where we cannot have a provable violation and instead
ned to make an accusation.

As an example lets take PO1: How can rule PO1 be violated?

PO1 is violated when a proposal for an old value V is sent by pi and we cannot
see a quorum of prevotes for that value in a previous round and then a
precommit of V or nil from pi in every intervening round until pi sends the
old proposal.

Unfortunately there is also a mistake in this rule due to not taking into
account line 55 of tendermint. If we take that into account we can't be sure
that we will see any precommits from pi. So the rule becomes weaker simply that
if pi proposes an old V then we should see a quorum of prevotes for V from a
prior round.

But we can strengthen it by changing the second term from "(nil v V)\*" to

"not a contiguous block of precommits sent by pi from the round where we saw
the prevote quorum threshold to the round preceding pi's old value proposal
where the last non nil precommit is for a value other than V"

That is quite a mouthful so maybe we could introduce some other notation.

I introduce ! to mean not and [x,y] to mean any combination of the comma
separated elements inside. Additionally I introduce + and * for use inside
square brackets to denote the cardinality of each element where * indicate zero
or more times and + indicates one or more times.and when used inside [..]
following an elemn

Old PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V ⇒  V` 
New PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  ![nil*,V'+] ⇒  V` 

So the new rule reads: if we see a proposal for an old value from pi then we
expect to see a quorum of prevotes in a prior round and a contiguous block of
precommits from pi which are not some combination of nil and V' with at least one instance of V'.
But that actually doesn't work since the rule at line 55 means we can't expect to see anything.

Try again:

I will introduce some more new notation:

I will use (..)+ to mean whatever in the brackets expanding to
fill all the rounds at least one.

we need to redefine the general rule for PO So the new rule says we
see prevotes in an old round then maybe see some nil or V precommits and then
don't see at least one precommit for V' followed by some possible precommits
for nil.

New PO: `M r'<r,PV ≺ (M r'≤r''<r,PC|pi)* ≺ !(M r''≤r'''<r,PC|pi)+ ≺ (M r'''≤r''''<r,PC|pi)* ⇒  Mr,P|pi`
New PO1:  `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v V ≺  V' ≺  nil ⇒  V` 

This notation is not great, its very verbose and doesn't cover the case where
we have the sequence "nil V' nil V'"

So new notation, I want to say that [nil\* V'+] means some sequence of V' and
nil with one or more V' and zero or more nil. The problem with this is it
doesn't fit with the approach taken to define a general rule and then provide
instances with values.

Ok I think that we actually need to introduce something else here the notation for no message at all.

So we can introduce ⊥ to mean no message, this saves us from having to use not (!) and the square brackets

Old PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V ⇒  V` 
New PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V v ⊥ ⇒  V` 

This still doesn't work though because it is still valid that pi could switch
to precommitting to V' and back to V by receiveing enough prevotes. But the
assumption is that we should see those prevotes since that is the first term of
PO1. So what we want to say is that from the prevotes for the highest round
with a quorum threshold for V we should not see a precommit for V'.

This also doesn't work, I think becuase it assumes that pi saw the latest prevotes and maybe they did not.

I think actually the only thing we can say for PO1 is that in some prior round there must have been a quorum of prevotes for V.

No we can say more, given that we assume we see all the messgaes in the network.

So what we want to say is that from the most recently observed quorum of
prevotes for V we should only witness precommits for V precommits for nil or no
precommit from pi in the intervening rounds till pi proposes V. It's difficult
to express this in a rule since we have no notation to mark the most recent of
something.

----------------------------------- I am here


# Rules classification

Some rules can always produce a concise proof of misbehaviour, others can only
produce accusations.

Here we classify rules into either group.

#f Concise proof


If we see the right side withouth the prerequisite left side or vice versa.

E.G

If we see a proposal for a new value V and on the left side a precommit for
anything other than nil we have misbehaviour.

or 

If we see on the left side precommits only for nil and then a proposal for an old value we also have misbehaviour. But actually the case is PN1 where it is assumed that the propose value is a new value.


PN1 (a prior precommit and the new value proposal)




---------------------Scrap section


# PO1 discussion
As an example lets take PO1: How can rule PO1 be violated?

PO1 is violated when a proposal for an old value V is sent by pi and we cannot
see a quorum of prevotes for that value in a previous round and then a
precommit of V or nil from pi in every intervening round until pi sends the
old proposal.

Unfortunately there is also a mistake in this rule due to not taking into
account line 55 of tendermint. If we take that into account we can't be sure
that we will see any precommits from pi. So the rule becomes weaker simply that
if pi proposes an old V then we should see a quorum of prevotes for V from a
prior round.

But we can strengthen it by changing the second term from "(nil v V)\*" to

"not a contiguous block of precommits sent by pi from the round where we saw
the prevote quorum threshold to the round preceding pi's old value proposal
where the last non nil precommit is for a value other than V"

That is quite a mouthful so maybe we could introduce some other notation.

I introduce ! to mean not and [x,y] to mean any combination of the comma
separated elements inside. Additionally I introduce + and * for use inside
square brackets to denote the cardinality of each element where * indicate zero
or more times and + indicates one or more times.and when used inside [..]
following an elemn

Old PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V ⇒  V` 
New PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  ![nil*,V'+] ⇒  V` 

So the new rule reads: if we see a proposal for an old value from pi then we
expect to see a quorum of prevotes in a prior round and a contiguous block of
precommits from pi which are not some combination of nil and V' with at least one instance of V'.
But that actually doesn't work since the rule at line 55 means we can't expect to see anything.

Try again:

I will introduce some more new notation:

I will use (..)+ to mean whatever in the brackets expanding to
fill all the rounds at least one.

we need to redefine the general rule for PO So the new rule says we
see prevotes in an old round then maybe see some nil or V precommits and then
don't see at least one precommit for V' followed by some possible precommits
for nil.

New PO: `M r'<r,PV ≺ (M r'≤r''<r,PC|pi)* ≺ !(M r''≤r'''<r,PC|pi)+ ≺ (M r'''≤r''''<r,PC|pi)* ⇒  Mr,P|pi`
New PO1:  `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v V ≺  V' ≺  nil ⇒  V` 

This notation is not great, its very verbose and doesn't cover the case where
we have the sequence "nil V' nil V'"

So new notation, I want to say that [nil\* V'+] means some sequence of V' and
nil with one or more V' and zero or more nil. The problem with this is it
doesn't fit with the approach taken to define a general rule and then provide
instances with values.

Ok I think that we actually need to introduce something else here the notation for no message at all.

So we can introduce ⊥ to mean no message, this saves us from having to use not (!) and the square brackets

Old PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V ⇒  V` 
New PO1: `#(Mr' ,PV|V) ≥ 2f + 1 ≺  nil v  V v ⊥ ⇒  V` 

This still doesn't work though because it is still valid that pi could switch
to precommitting to V' and back to V by receiveing enough prevotes. But the
assumption is that we should see those prevotes since that is the first term of
PO1. So what we want to say is that from the prevotes for the highest round
with a quorum threshold for V we should not see a precommit for V'.

This also doesn't work, I think becuase it assumes that pi saw the latest prevotes and maybe they did not.

I think actually the only thing we can say for PO1 is that in some prior round there must have been a quorum of prevotes for V.

No we can say more, given that we assume we see all the messgaes in the network.

So what we want to say is that from the most recently observed quorum of
prevotes for V we should only witness precommits for V precommits for nil or no
precommit from pi in the intervening rounds till pi proposes V. It's difficult
to express this in a rule since we have no notation to mark the most recent of
something.
