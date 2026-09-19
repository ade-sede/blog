# Reasoning about software architecture

Across 70 years of software engineering, our industry has had all kinds of debates about which software architecture is best.
We argue semi-regularly about how to organise and reason about code. About which way would be most modular, which would be most maintenable, which would be easiest to deploy, etc... We even argue about which of these properties we should optimise for! We reinvent old ideas every few years and claim it will change how we build software forever.

However, if we take a step back we can identify a common pattern: most software is built upon the concept of _layers_.

## Layers as our primitive building block

In software architecture, a layer represents a logical grouping. Anything goes: a grouping of concepts, of responsibilities, of related operations, etc...
A layer is defined by said grouping, and by its interactions with other layers. Using layers, we are provided with a mental model and tools to help us reason about systems.

The first tool and most valuable tools provided by layers is abstraction.
For example, when we write a web application we tend to focus on _what to display_ and don't give much thoughts to _how it is displayed_. The _browser layer_ offers us a contract: if we provide correct HTML and CSS, the browser will somehow be able to render it.
Some people will tell you the how and the what cannot be entirely dissociated and impact each other in many ways[^1], that the best craftsmen must master the _how_ if they want to produce high quality software, and I'd say they're right! But even if all abstractions are imperfect, they can still be valuable.
Us humans are bounded by our cognitive abilities and any opportunity to isolate problems and consider them one at a time is a win.

The second tool is modularity.
If we accept interactions between layers can be defined by their contracts, it stands to reason two different layers offering the same exact contract can be used interchangeably. Another way to look at it is defining the contract is a way to express our requirements for a layer we depend on.
For example, any browser will be able to render our web application as long as it implements the same standard for HTML & CSS.

A third is hierarchy

layering is also what decides failures modes
its possible for me to reason about them in isolation - everywhere we made this isolation impossible, we created what is called coupling.
being able to think one at a time doesn't mean knowledte of both sides is not required or that ignorance is good

picking constraints = picking itneraction and contract = picking what we are coupled to

Modularity is about boundaries (horizontal separation), while Layering is about hierarchy and dependency direction (vertical separation).

forcing to care = leak
abstraction can hurt you, sometimes:

- as it solves constraints it can create new ones
- can givew a false sense of security
  how yhou craft the contract has big impact

defining the cont ract and evaluating two deps actually rsepect the same exact contract without leaking is not so easy: example, problem compat issues

defining an architect
being fluent in describing responsibilities and interactions
being fluent in communicating that picture: code as artefact, code as shared language, ubigituous language
an architect is a master layerist, picks the right axioms given a set of constraints and known unknowns

where to put boundary, why.
The act of architecting = taking clear decisions with our eyes wide open
with the knowledge of all constrraints as well as all known unknowns, which shape should the software take?
the architect is the person who understands why it's the right shape for the building and

comp sci, craftmanship, product sense, business sense, everything collapsing at once. equilibriste
melting pot of evertything.
software architecture as a discipline is developing a methodology to indentify and work through the mess

what is an astronaut architect? is the act of thinking a crime?

Layers are how humans cope with complexity → the architect's job is picking the right layer boundaries → this requires knowing your constraints and unknowns → the hard part isn't thinking, it's communicating that picture to others → here's where most teams go wrong (infra leaking into domain?)

## The role of an architect

## Teleporting ideas: from the architect's mind, to reality, to other

High BW communication
ubiquitous languague

sum of constraints

[^1]: example: performance
