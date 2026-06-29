"""
Treina um classificador de URLs suspeitas (phishing/malware).
Usa features estruturais da URL — nao precisa baixar a pagina.
"""
import re
import pickle
from urllib.parse import urlparse
import numpy as np
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report


def extract_features(url: str) -> list:
    """Extrai features estruturais de uma URL."""
    parsed = urlparse(url if "://" in url else "http://" + url)
    host = parsed.netloc
    path = parsed.path

    return [
        len(url),                                    # tamanho total
        len(host),                                   # tamanho do host
        url.count("."),                              # numero de pontos
        url.count("-"),                              # hifens (comum em phishing)
        url.count("@"),                              # @ (redirecionamento suspeito)
        url.count("/"),                              # profundidade do path
        1 if re.search(r"\d+\.\d+\.\d+\.\d+", host) else 0,  # IP no lugar de dominio
        1 if parsed.scheme == "https" else 0,        # tem https
        len(re.findall(r"\d", url)),                 # quantidade de digitos
        1 if any(w in url.lower() for w in           # palavras-isca
            ["login", "secure", "account", "verify", "update", "bank"]) else 0,
        len(host.split(".")),                        # numero de subdominios
    ]


# Dataset sintetico: URLs legitimas vs suspeitas
legit = [
    "https://google.com", "https://github.com/user/repo",
    "https://www.amazon.com/products", "https://wikipedia.org/wiki/Python",
    "https://news.ycombinator.com", "https://stackoverflow.com/questions",
    "https://www.microsoft.com", "https://docs.python.org/3/",
    "https://www.bbc.com/news", "https://twitter.com/home",
    "https://www.netflix.com", "https://linkedin.com/in/profile",
    "https://reddit.com/r/programming", "https://medium.com/article",
    "https://www.apple.com/iphone", "https://aws.amazon.com/ec2",
    "https://kubernetes.io/docs", "https://golang.org/doc",
    "https://www.spotify.com", "https://www.dropbox.com",
]

phishing = [
    "http://paypa1-login.secure-verify.xyz/account",
    "http://192.168.1.1/bank-login/verify.php",
    "http://secure-update-account.ru/login?id=12345",
    "http://amaz0n-account-verify.tk/secure",
    "http://login-microsoft-verify.ml/update-account",
    "http://apple-id-locked-verify.cf/unlock",
    "http://bank-of-america-secure.ga/login.html",
    "http://verify-your-account-now.xyz/paypal",
    "http://192.0.2.55/secure-login-bank-verify",
    "http://netflix-billing-update.tk/account/verify",
    "http://signin-ebay-secure-verify.ml/login",
    "http://update-payment-info-secure.ga/bank",
    "http://confirm-account-suspended.cf/verify-now",
    "http://google-security-alert-verify.xyz/login",
    "http://dropbox-shared-file-login.tk/secure",
    "http://instagram-copyright-verify.ml/appeal",
    "http://whatsapp-premium-verify.ga/account",
    "http://10.0.0.1/admin-login-secure-bank",
    "http://facebook-security-checkpoint.cf/verify",
    "http://your-package-delivery-verify.xyz/track",
]

# monta X e y
urls = legit + phishing
labels = [0] * len(legit) + [1] * len(phishing)  # 0=ok, 1=suspeito

X = np.array([extract_features(u) for u in urls])
y = np.array(labels)

X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.25, random_state=42, stratify=y
)

clf = RandomForestClassifier(n_estimators=100, random_state=42)
clf.fit(X_train, y_train)

print("=== Avaliacao ===")
print(classification_report(y_test, clf.predict(X_test),
      target_names=["legitima", "suspeita"]))

with open("model/url_classifier.pkl", "wb") as f:
    pickle.dump(clf, f)

print("Modelo salvo em model/url_classifier.pkl")
