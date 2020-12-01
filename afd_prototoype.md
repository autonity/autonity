# How can rules be violated?

we use pi to define the participant under scrutiny. And we understand the m I
used to enjoy those things more than my son sometimeeaning of * when applied to
an expression in brackets to mean that that expression expands to fill the
range of rounds defined in the brackets but that that range of rounds may be
empty.

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
precommits from pi which are not some combination of nil and v' with at least one precommit for V'.
But that actually doesn't work since
the rule at line 55 means we can't expect to see anything.

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

