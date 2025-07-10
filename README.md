```
██████╗ ██╗  ██╗ █████╗ ██╗██╗     
██╔══██╗██║  ██║██╔══██╗██║██║     
██████╔╝███████║███████║██║██║     
██╔═══╝ ██╔══██║██╔══██║██║██║     
██║     ██║  ██║██║  ██║██║███████╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚══════╝

P o o r - m a n ’ s   H i g h l y  
  A v a i l a b l e   I n t e r n e t  
     L o a d  B a l a n c e r
```
![PHAIL Status](https://img.shields.io/badge/PHAIL-PASSING-green?style=flat-square&logo=cloudflare)
![Budget](https://img.shields.io/badge/Budget-$5%2Fmo-orange?style=flat-square&logo=money)
![Uptime](https://img.shields.io/badge/Uptime-99.9%25-yellow?style=flat-square)
![Vibes](https://img.shields.io/badge/Vibes-Unhinged-ff69b4?style=flat-square)
![Failover](https://img.shields.io/badge/Failover-Sorta_Working-blueviolet?style=flat-square)

---

##### When you wanna vibe with “99.9% uptime” but your wallet says “nah.”


## ✨ wtf is PHAIL?

PHAIL is your **$5 “global load balancer”** knockoff — a squad of feral Bun goblins that:
- Gossip about which nodes are alive 💀
- Elect a ✨ main character ✨ to flip Cloudflare DNS
- Rage-flips A records when your API decides to head out.
- Throws hands internally if the leader dies

PHAIL is your broke bestie for pretending you run multi-cloud failover.

💸 No GSLB.

💀 No fancy Anycast.

🤠 Just vibes, health checks, and DNS record flips.

## 🤡 how it works

1️⃣ Deploy 3 PHAIL goblins across some crusty $5 VPSs. (Or just 1 if you don't want HA)

2️⃣ Each goblin vibe checks your servers 24/7

3️⃣ Leader goblin hits Cloudflare API to yeet dead nodes from your A records.

4️⃣ If the leader goblin trips & dies, another goblin shouts “Bet, my turn.”

## 💀 why even use PHAIL?

• It’s cheaper than any other load balancer.

• Route53 wants your money, PHAIL wants your love

• Real load balancers are for corporate sellouts

• You want to say “we’re multi-cloud” unironically.

## 🗿 features (kinda)

✅ Multi-cloud, multi-VPS, works on everything

✅ Cloudflare DNS flips when your node ghosts

✅ Absolutely no Kubernetes — you’re better than that

✅ Can run on a Nokia 3310 (probably)

---

## 🗃️ Architecture
```
           [ Cloudflare DNS ]
                    🫧
┌────────────────────┴───────────────┐
|                    |               |
[ DigitalOcean ][ Hetzner ][ Basement Server ]
|  🫃 Goblin    🫃 Goblin       🫃 Goblin
|                    |               |
```

## 🧃 quickstart (development)

1️⃣ git clone this scuffed repo

2️⃣ bun install

3️⃣ bun run dev

## 😭 known issues
• DNS TTL isn’t magic, you may encounter up to 1 min of downtime per failure.

## 🫡 license

GPLv3. No refunds. No uptime SLA. If you PHAIL, that’s on you.

## 🫶 PHAIL: You only live once. Keep your servers up. Mostly.

## 👀 made with:
• Bun

• bad life choices

• a burning desire not to pay AWS