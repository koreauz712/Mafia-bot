import random
from telegram import Update
from telegram.ext import (
    ApplicationBuilder,
    CommandHandler,
    ContextTypes
)

TOKEN = "8260077131:AAE66L2kNNoMAvQCDh3cfNd2X8KV9lfifNM"
games = {}

roles_list = [
    "🕵 Mafia",
    "👮 Sheriff",
    "💉 Doctor",
]

# 🎮 CREATE
async def create(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id

    games[chat_id] = {
        "players": [],
        "started": False
    }

    await update.message.reply_text(
        "🎮 Mafia o‘yini yaratildi!\n\n"
        "Qo‘shilish uchun /join bosing"
    )

# 👥 JOIN
async def join(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id
    user = update.effective_user

    if chat_id not in games:
        await update.message.reply_text("❌ Avval /create")
        return

    players = games[chat_id]["players"]

    if user.id in [p["id"] for p in players]:
        await update.message.reply_text("⚠️ Siz allaqachon qo‘shilgansiz")
        return

    players.append({
        "id": user.id,
        "name": user.first_name
    })

    await update.message.reply_text(
        f"✅ {user.first_name} o‘yinga qo‘shildi\n"
        f"👥 Odamlar: {len(players)}"
    )

# 🎭 START GAME
async def startgame(update: Update, context: ContextTypes.DEFAULT_TYPE):
    chat_id = update.effective_chat.id

    if chat_id not in games:
        return

    players = games[chat_id]["players"]

    if len(players) < 3:
        await update.message.reply_text(
            "❌ Kamida 3 ta odam kerak"
        )
        return

    # roles
    roles = roles_list.copy()

    while len(roles) < len(players):
        roles.append("👨 Civil")

    random.shuffle(roles)

    # PRIVATE ROLE SEND
    for i, player in enumerate(players):
        try:
            await context.bot.send_message(
                chat_id=player["id"],
                text=f"🎭 Sizning rolingiz:\n\n{roles[i]}"
            )
        except:
            pass

    await update.message.reply_text(
        "🔥 O‘yin boshlandi!\n"
        "📩 Rollar private chatga yuborildi"
    )

# 🤖 APP
app = ApplicationBuilder().token(TOKEN).build()

app.add_handler(CommandHandler("create", create))
app.add_handler(CommandHandler("join", join))
app.add_handler(CommandHandler("startgame", startgame))

print("🔥 Mafia bot ishlayapti")
app.run_polling()
