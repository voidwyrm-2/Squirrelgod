#! usr/bin/env python3
from pathlib import Path


def inputr(msg: str, nonempty: bool = False) -> str:
    if not nonempty:
        msg += " (optional, leave empty to skip)"
    while True:
        i = input(msg + "\n").strip()
        if i == "" and nonempty:
            continue
        return i


print("creating required files...")

config = Path("config.txt")
offerings = Path("offeringCount.txt")

print(f"creating '{config}'...")

with open(config, "wt" if config.exists() else "xt") as cf:
    token = inputr("what is your bot token?", True)
    invite_link = inputr("what is your bot invite link?")
    source_link = inputr("what is the link to your repo?")
    online_annoucement_channel = inputr(
        "what is the ID of the channel for sending online announcements?"
    )
    cf.write(f"""token:string: {token};
invite_link:string: {invite_link};
source_link:string: {source_link};
announce_channel:string: {online_annoucement_channel};
online_messages:string:;""")

print(f"created '{config}'")

print(f"creating '{offerings}'...")

if offerings.exists():
    with open(offerings, "w+rt") as of:
        content = of.read()
        try:
            int(content)
        except Exception:
            of.write("0")
else:
    with open(offerings, "wt") as of:
        of.write("0")

print(f"created '{offerings}'")
