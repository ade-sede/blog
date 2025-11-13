# The tyranny of simple

When building software products, a concept frequently brought up is _simplicity_. The idea that the way we approach any and all topics should be uncomplicated, straightforward.
Simplicity has become a mantra, initially chanted by some very successful people in the software industry and later echoed by anyone trying to imitate their success.
Everyone agrees simple is best. No one sets out to build something overly complex on purpose.
Yet it keeps coming up: _"This is over-engineered"_, _"I don't know, looks complicated"_.
I believe the root cause comes from a fundamental misalignment: although modern software organisations expect simplicity, few of them actually spend time establishing what _simple_ means.
Without shared understanding it is an obsession for an arbitrary standard no one can meet reliably.
If we truly believe simplicity is something to strive for, we need to develop the necessary vocabulary and culture for productive conversations to occur.
Invoking the _it's too complex_ argument should never be a conversation ender but rather a signal we need to start talking and find out where our mutual expectations started diverging.

In this essay I want to explore what simplicity is and how organisations can create a culture where it is more than a buzzword.

## Attempting a definition of simplicity

Counter-intuitively, simplicity is not so easy to define. Mainly because it spans several dimensions[^1]:

Simplicity is the opposite of complexity.
Simplicity is not a binary property.
Simplicity is a tool: the simpler something is, the easier it becomes to evolve, verify, debug, and reason about.
Something simple does not drown us under its _cognitive load_.
Something simple is _intuitive_ and _self contained_: very little specialized knowledge is required and it can be understood simply by looking at it a few times.
Something simple is _predictable_: once we understand it we are able to predict how it will behave in different scenarios.
Simple things _scale_ and are _composable_: they can be assembled together and remain simple. In contrast, when assembling complex things the complexity compounds and explodes, overloading our cognitive abilities.
Simplicity is a property any concept or artifact[^2] can have.
It is also important to make the distinction between _essential complexity_, unavoidable because it stems directly from the problem we are working on and _incidental complexity_, which is entirely avoidable and always self-inflicted.

Every bit of nuance in this definition is a gap where misalignment can grow.

The first gap I see is the understanding of the problem we are solving. It is possible two contributors do not see the same amount of essential complexity.
The most obvious example being contributor A designing to make enough room for future requirements while contributor B targets the minimum implementation for current requirements.
Maybe there is a justifiable need for a future proof design, because there we foresee an immediate, plausible need to scale. Maybe it doesn't make the list because we need to move fast and validate some hypothesis first.
In any case, we can easily understand why contributor B would see contributor A's solution as over-engineered and adding incidental complexity: it is literally engineered for a different, wider set of requirements.

The second gap is the personal bias of each individual. Perception is subjective and we do not realise how much of it is shaped by our past experiences.
Everyone's threshold for cognitive load is different. And most importantly: simplicity is easily confused with familiarity.
Once we have already invested time and effort into understanding something it feels intuitive, independently of how costly the initial investment was.

Simplicity is much harder than most people expect, and considering how subtle the nuances are I do not think it is reasonable to expect every individual to share the same exact perception of simplicity.
This alignment is something that must be crafted intentionally and weaved into the culture of the organisation.

## How to create the right culture

Building the culture of an organisation is as much about stating convictions as it is about embodying them. If either aspect is missing the culture does not take hold and is quickly forgotten.

### State convictions out loud

The first thing any org striving for simplicity should do is to publish a manifesto clearly laying out their view of simplicity. It's important for this manifesto to go beyond guidelines and get into why this view emerged in the first place.
Performant orgs need contributors to be able to make their own calls, and they can only do so if they have a deep understanding of the drivers behind the philosophy.
The manifesto should be updated every time the org crosses a major milestone. With new challenges come a different set of requirements and it is only natural for the philosophy to evolve.

### The contributors protect the bar

With the convictions now in the open for everyone to see, the next priority is ensuring the entire organisation lives up to them. I believe the only sustainable way to achieve this is to distribute the responsibility amongst all contributors. Everyone must feel empowered to speak up whenever they feel a change does not meet the bar, or the bar is being misused.

To make this work in practice, the organisation needs to establish a simple but powerful norm: whoever invokes simplicity as an argument bears the burden of proof. This is a forcing function, a structural mechanism baked into the culture. Saying "this is too complex" or "my approach is simpler" is not the end of a discussion, it's the beginning.

The person making the claim must articulate specifically what makes something complex or simple, referencing the shared principles. Is it the cognitive load? The amount of specialized knowledge required? The unpredictability of edge cases? The way complexity compounds when combined with other parts of the system? Without this burden of proof, "too complex" remains a gut feeling that people genuinely experience but cannot examine. By requiring concrete reasoning, it helps contributors distinguish between legitimate complexity concerns and the discomfort of unfamiliarity. The goal is to move beyond instinct to articulate what specific property of simplicity is being violated.

This might sound overly formal, but there's a good reason for it: conversations about simplicity can easily become personal.

When someone presents a solution they believe is elegant, they've often invested significant mental effort into crafting it. Software engineers, perhaps more than other professionals, can conflate the cleverness of their solutions with their own intelligence and value. An abstraction or pattern they've designed becomes a reflection of their capabilities. When someone pushes back with "this is too complex", it can land as "you over-engineered this" or worse, "you're trying to be too clever" rather than as technical feedback. The person giving feedback might genuinely be concerned about maintainability, but the person receiving it can hear an attack on their competence.

Without a shared vocabulary and explicit principles to reference, these discussions devolve into subjective taste arguments where both parties feel personally invested in being "right". The manifesto and clear vocabulary serve as a neutral third party. They allow contributors to debate the merits of a solution against established principles rather than against each other.

Building this culture means making these structured, constructive conversations the norm. It's not enough to simply tell people "pushback is allowed". The organisation needs to model what productive pushback looks like. So what does good pushback actually look like? When someone feels a solution is too complex, they should ask themselves a series of clarifying questions. Is the complexity in the artifact itself (verbose code, many moving parts, unclear flow) or is it in the underlying concept? This maps back to essential versus incidental complexity. If the concept itself is inherently complex because of the problem domain, then some amount of implementation complexity might be unavoidable. The question becomes: is this the minimum complexity required, or are we adding incidental complexity on top?

Another useful lens: what piece of knowledge would make this feel simple? If the answer is "I'd need to understand X pattern" or "I'd need to learn Y library", that's not necessarily a mark against the solution. It's identifying a knowledge gap. The question then becomes whether that knowledge investment is justified by the problem we're solving.

Finally, when advocating for a different approach based on their own experience, contributors should be explicit about the connection: which aspects of their past experience actually apply here, and which aspects of this problem are genuinely different?

### Resilience to change

Essential complexity grows as organizations mature. Early stage teams optimize for speed, but as products scale new constraints emerge: uptime requirements, data consistency, regulatory compliance, etc... The problems become more complex and solutions must grow with them.

The challenge is knowing which complexity to accept and which to reject. This is where the burden of proof and shared vocabulary become essential. By articulating which property of simplicity is being violated and referencing the manifesto, teams can separate incidental from essential complexity.

Without this discipline, familiarity bias creeps back in: teams reject unfamiliar but necessary complexity while accepting familiar but avoidable complexity. The manifesto must evolve as the organization crosses major milestones. What worked for an early-stage startup won't necessarily work at scale.

To make this manageable, organizations must invest in building pits of success: structural mechanisms that make coping with new, increased requirements feel as simple as the previous ones. The goal is to lower the floor and raise the bar simultaneously, making it clearer what's expected while expanding what the team can handle.

One example of such a pit of success is curating two lists alongside the manifesto: necessary domain knowledge and active constraints. The knowledge list captures patterns, libraries, and domain concepts contributors need to be effective. The constraints list documents non-negotiables: performance targets, uptime requirements, regulatory obligations, etc... Both lists should be maintained with the same rigor as the manifesto. This transparency transforms subjective debates into concrete tradeoff discussions.

Diverse experiences become invaluable here. Contributors who have worked at different scales and in different domains bring varied mental models that help the entire organization recognize when complexity is justified and when it's merely unfamiliar. Valuing diverse experiences builds resilience against familiarity bias while increasing the threshold for essential complexity.

## Conclusion

Organizations must define what simplicity means for their context. They should build forcing functions into their culture to ensure simplicity always comes on top. They should maintain explicit lists of necessary knowledge and constraints. They should value diverse experiences. These aren't revolutionary ideas, but they require consistent effort. Writing and maintaining a manifesto takes time. Articulating why something feels wrong instead of relying on instinct is uncomfortable. Recognizing that familiarity bias might be misleading requires humility.

When organizations do this work well, debates about complexity become learning opportunities. Contributors develop shared intuition. Teams move faster because they're aligned on principles, not arguing about taste. And simplicity becomes what it should be: a tool that accelerates teams and helps them scale.

[^1]: While writing this article I stumbled upon the excellent talk ["Simple Made Easy"](https://www.youtube.com/watch?v=SxdOUGdseq4) by Rich Hickey. I don't agree with all of the takeaways, but it is definitely the best explanation of what simplicity is I have ever heard. So much so I feel stupid writing this section of the article... Highly recommend watching.

[^2]: Examples: an idea, a system, a piece of code, a piece of documentation, or the experience of a user.
