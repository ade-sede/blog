# The tyranny of 'simple'

_Simple is best_.
This is a pretty uncontroversial statement in modern software teams... Everyone agrees. So why does the topic keep coming up?
Engineer A pushes some code for review. Engineer B doesn't like it.

- "It's over-engineered! We should just do XYZ.", says engineer B.
- "I don't know what you mean, this really isn't that complex.", objects engineer A.
- "Stop arguing, it doesn't matter.", concludes engineer C.

Are any of them right? Is there a point to figuring out the answer? What does it mean for something to be simple?

## A crude definition

Counter-intuitively, defining simplicity is not easy. Though that won't stop me from trying!

Simplicity is the opposite of complexity.
Simplicity is not a binary property.
Simplicity is a tool: the simpler something is the easier it is to verify, reason about and modify.
Something simple does not drown us under its _cognitive load_.
Something simple is intuitive and self contained: very little specialized knowledge is required and it can be understood simply by looking at it a few times.
Something simple is predictable: once we understand it we are able to predict how it will behave in different scenarios.
Simple things scale and are composable: they can be assembled together and remain simple. In contrast, when assembling complex things the complexity compounds and explodes, overloading our cognitive abilities.
Simplicity is a property any concept or artifact[^1] can have.
There are several kinds of complexity. _Essential complexity_ is unavoidable because it stems directly from the problem we are working on. _Incidental complexity_ is entirely avoidable and always self-inflicted.

Taking a step back, a lot of it has to do with our ability to construct a coherent mental model and navigate it.
Let's try to exercise this definition.

One of the most relatable sources of accidental complexity in software comes from misplacing responsibilities.
When responsibilities are misplaced we accidentally create coupling between systems that should be entirely separated. And when two systems are coupled, we can't reason about one without reasoning about the other, which naturally increases the cognitive load.

Another example is bad abstraction.
Abstraction is a powerful tool: done right, it reduces how much you need to know and understand about the underlying system. Done wrong, it can pin you down under specific assumptions and you become unable to make certain changes.

Now that we have a definition, flawed as it may be, does it mean we are able to say for sure who of engineer A and engineer B is right?

## The trap of familiarity

Our perception of what is simple is extremely subjective.
We are not all necessarily equal when it comes to building mental models: we do not have access to the same pool of knowledge, we do not have the same threshold to define when the cognitive load is overwhelming, we do not have access to the same techniques for assembling concepts and navigating relationships.
Our past experiences color our perception, and this is why most often simplicity is confused with familiarity.
While it is possible for someone to have a genuinely bad idea, ultimately engineer A and engineer B are arguing about having a different mental model, fueled by different experiences.

So is engineer C right? Maybe none of this matters and trying to reach a result universally perceived as simple is impossible and not worth pursuing... I refuse to believe this is the case.

## Turning the disagreement into an opportunity

It is widely accepted software development is a team sport. Many individuals join forces and collaborate to deliver a working product, and then proceed to maintain it over a long period of time.
Simplicity has a direct impact on our ability to reason about and modify our software and the upside is considerable. Not only do we want the current team to feel at ease but we also want future team members to be productive and in control.

If the main blocker to building a shared mental model are our diverging experiences and knowledge pool, all we need is to make them converge.
A shared knowledge pool is something that can be built from scratch: documentation, examples, workshops.
Past experiences can be boiled down to important learnings and shared publicly: retrospectives, pairings.
Once we unravel 'why' there is a disagreement, we can identify all the gaps in the mental models, debate them on their merits and make decisions based on which are relevant to the task at hand.

# Conclusion

Leaders dictating simplicity without arming their teams to actually reach it are creating an unfair standard. It is in every organisation's interest to limit incidental complexity, yet very few of them actually define what it looks like and make the mistake of assuming it is a universal concept.
Where young organisations might see debates and disagreements as unproductive whining and shut them down, mature organisations will recognise them as a signal a shared mental model is missing.
The difference between both is whether or not they have built the culture and shared vocabulary necessary to defuse the trap of familiarity and boil objections down their core principles, instead of a dead-end "he said/she said" situation.

[^1]: Examples: an idea, a system, a piece of code, a piece of documentation, or the experience of a user.
